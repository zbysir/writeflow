package writeflow

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
	"strconv"
	"strings"
)

type WriteFlow struct {
	cmds map[string]CMDer
}

func NewShelFlow() *WriteFlow {
	return &WriteFlow{
		cmds: map[string]CMDer{},
	}
}

func (f *WriteFlow) RegisterCmd(taskName string, fun CMDer) {
	f.cmds[taskName] = fun
}

func execFunc(ctx context.Context, fun interface{}, params []interface{}) (rsp []interface{}, err error) {
	callParams := []reflect.Value{reflect.ValueOf(ctx)}

	for _, p := range params {
		callParams = append(callParams, reflect.ValueOf(p))
	}
	funv := reflect.ValueOf(fun)
	ty := funv.Type().NumIn()
	for i := 0; i < ty; i++ {
		wantp := funv.Type().In(i)
		inp := callParams[i].Type()

		//fmt.Printf("wantp:%v, inp:%v %v\n", wantp.String(), inp.String(), inp.AssignableTo(wantp))

		// TODO 如果目标是数组，则使用 Append 而不是直接赋值，来源可以支持数组 Item

		// 如果类型不匹配，则尝试通过 json 转换
		if !inp.AssignableTo(wantp) {
			bs, _ := json.Marshal(callParams[i].Interface())
			w := reflect.New(wantp)
			err = json.Unmarshal(bs, w.Interface())
			if err != nil {
				return nil, fmt.Errorf("can not convert %v to %v, err: %w", inp.String(), wantp.String(), err)
			}

			callParams[i] = w.Elem()
		}

		//if wantp.String() == "[]string" && inp.String() == "[]interface {}" {
		//	callParams[i] = reflect.ValueOf(interfaceTo[string](callParams[i].Interface().([]interface{})))
		//}

	}

	rv := funv.Call(callParams)

	var rerr error
	l := len(rv)
	for i, v := range rv {
		if i == l-1 {
			last := v
			switch last.Kind() {
			case reflect.Interface:
				err, ok := last.Interface().(error)
				if ok {
					rerr = err
					continue
				}
			}

			rsp = append(rsp, v.Interface())
		} else {
			rsp = append(rsp, v.Interface())
		}
	}

	return rsp, rerr
}

// 所有的依赖可以并行计算。
// 这是通过代码逻辑不好描述的
//
// appendName-1:
//   cmd: appendName
//   input:
//     - _args[0]
//     - _args[1]
// hello:
//   cmd: hello
//   input:
//     - appendName-1[0]
//
// END:
//   input:
//     - hello[0]

// flow: 流程定义
// job: flow 由 多个 job 组成
// cmd: job 可以调用 cmd

type JobInput struct {
	// _args[0]
	JobName   string
	RespIndex int
	// {a: _args[1]}
	Object map[string]JobInput
}

type JobDef struct {
	Name   string
	Cmd    string
	Inputs []JobInput
}
type FlowDef struct {
	Jobs map[string]JobDef
}

// SpanInterface 特殊语法，返回值
type SpanInterface []interface{}

type YFlow struct {
	Version string          `yaml:"version"`
	Flow    map[string]YJob `yaml:"flow"`
}

type YJob struct {
	Cmd     string        `yaml:"cmd"`
	Inputs  []interface{} `yaml:"inputs"`
	Depends []string      `yaml:"depends"`
}

func (f *YFlow) ToFlowDef() FlowDef {
	jobs := map[string]JobDef{}
	for name, v := range f.Flow {
		jobs[name] = v.ToJobDef(name)
	}

	return FlowDef{Jobs: jobs}
}

func (j *YJob) ToJobDef(name string) JobDef {
	var inputs []JobInput
	for _, item := range j.Inputs {
		switch item := item.(type) {
		case string:
			// _args[1]
			ss := strings.Split(item, "[")
			taskName := ""
			var respIndex int64
			if len(ss) == 2 {
				taskName = ss[0]
				respIndex, _ = strconv.ParseInt(ss[1][0:len(ss[1])-1], 10, 64)
			} else {
				taskName = ss[0]
				respIndex = -1 // -1 表示就当成数值传递
			}

			inputs = append(inputs, JobInput{
				JobName:   taskName,
				RespIndex: int(respIndex),
				Object:    nil,
			})
		case map[string]interface{}:
			// {name: args[0]}
			// TODO object
		}
	}

	return JobDef{
		Name:   name,
		Cmd:    j.Cmd,
		Inputs: inputs,
	}
}

func (f *WriteFlow) parseFlow(flow string) (FlowDef, error) {
	var flowDefI YFlow
	err := yaml.Unmarshal([]byte(flow), &flowDefI)
	if err != nil {
		return FlowDef{}, fmt.Errorf("unmarshal flow err: %v", err)
	}
	def := flowDefI.ToFlowDef()

	return def, nil
}

func (f *WriteFlow) ExecFlow(ctx context.Context, flow string, params []interface{}) (rsp []interface{}, err error) {
	f.RegisterCmd("_args", FunCMD(func(ctx context.Context) SpanInterface {
		return params
	}))
	def, err := f.parseFlow(flow)
	if err != nil {
		return nil, err
	}

	fr := FlowRun{
		flowDef: def,
		cmdRsp:  map[string][]interface{}{},
		//args:    params,
		cmd: f.cmds,
	}
	rsp, err = fr.ExecJob(ctx, "END")
	if err != nil {
		return
	}

	return
}

type FlowRun struct {
	cmd     map[string]CMDer
	flowDef FlowDef
	cmdRsp  map[string][]interface{}
	//args    []interface{}
}

func (f *FlowRun) ExecJob(ctx context.Context, jobName string) (rsp []interface{}, err error) {
	jobDef := f.flowDef.Jobs[jobName]
	inputs := jobDef.Inputs

	//log.Printf("exec: %s, inputs: %v", jobName, inputs)
	dependValue := []interface{}{}
	for _, i := range inputs {
		var rsp interface{}
		if f.cmdRsp[i.JobName] != nil {
			//log.Printf("i.JobName %v: %+v", i.JobName, f.cmdRsp[i.JobName])
			// cache
			if i.RespIndex == -1 {
				rsp = f.cmdRsp[i.JobName]
			} else {
				rsp = f.cmdRsp[i.JobName][i.RespIndex]
			}
		} else {
			rsps, err := f.ExecJob(ctx, i.JobName)
			if err != nil {
				return nil, fmt.Errorf("exec task '%s' err: %v", i.JobName, err)
			}
			if len(rsps) == 1 {
				// 特殊语法，展开第一个元素
				switch rsps[0].(type) {
				case SpanInterface:
					rsps = rsps[0].(SpanInterface)
				}
			}

			//log.Printf("i.rsps %v: %+v", i.JobName, rsps)
			if i.RespIndex == -1 {
				rsp = rsps
			} else {
				rsp = rsps[i.RespIndex]
			}
			f.cmdRsp[i.JobName] = rsps
		}

		dependValue = append(dependValue, rsp)
	}

	cmd := jobDef.Cmd
	if cmd == "" {
		cmd = jobName
	}
	c, ok := f.cmd[cmd]
	if ok {
		rsp, err := execFunc(ctx, c, dependValue)
		if err != nil {
			return nil, fmt.Errorf("exec task '%s' err: %w", jobName, err)
		}

		return rsp, err
	} else {
		return dependValue, nil
	}

}

type Task interface {
	Do(ctx context.Context, req []interface{}) (rsp []interface{}, err error)
}
