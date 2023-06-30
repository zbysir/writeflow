package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/internal/pkg/ws"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/pkg/modules/builtin"
	"github.com/zbysir/writeflow/pkg/writeflow"
	"time"
)

type Flow struct {
	flowRepo     repo.Flow
	sysRepo      repo.System
	wirteflow    *writeflow.WriteFlow
	ws           *ws.WsHub
	pluginStatus []PluginStatus
}
type PluginStatus struct {
	model.PluginSource
	Status string
	Error  string
}

func NewFlow(flowRepo repo.Flow, sysRepo repo.System) (*Flow, error) {
	f := &Flow{
		flowRepo:  flowRepo,
		wirteflow: nil,
		sysRepo:   sysRepo,
		ws:        ws.NewHub(),
	}

	err := f.ReloadWriteFlow(context.Background())
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReloadWriteFlow 如果有插件更改，需要重新加载
func (u *Flow) ReloadWriteFlow(ctx context.Context) error {
	wf := writeflow.NewWriteFlow()

	wf.RegisterModule(builtin.New())

	setting, err := u.sysRepo.GetSetting(ctx)
	if err != nil {
		return err
	}

	// 加载插件
	u.pluginStatus = []PluginStatus{}
	pm := writeflow.NewGoPkgPluginManager(nil)

	// 设置一个默认的插件
	if len(setting.Plugins) == 0 {
		setting.Plugins = []model.PluginSource{
			{
				Url:    "https://github.com/zbysir/writeflow-plugin-llm",
				Enable: true,
			},
		}
	}
	for _, p := range setting.Plugins {
		if p.Enable == false {
			continue
		}

		plug, err := pm.Load(p.Url)
		if err != nil {
			u.pluginStatus = append(u.pluginStatus, PluginStatus{
				PluginSource: p,
				Status:       "error",
				Error:        err.Error(),
			})
		} else {
			err = plug.Register(wf)
			if err != nil {
				u.pluginStatus = append(u.pluginStatus, PluginStatus{
					PluginSource: p,
					Status:       "error",
					Error:        err.Error(),
				})
			} else {
				u.pluginStatus = append(u.pluginStatus, PluginStatus{
					PluginSource: p,
					Status:       "enable",
					Error:        "",
				})
			}
		}
	}

	u.wirteflow = wf

	return nil
}

func (u *Flow) GetComponents(ctx context.Context) (cs []writeflow.CategoryWithComponent, err error) {
	cs = u.wirteflow.GetComponentList()
	return cs, nil
}

func (u *Flow) GetComponentByKey(ctx context.Context, key string) (cs writeflow.Component, exist bool, err error) {
	cs, exist, err = u.wirteflow.GetComponentByKey(key)
	return
}

type RunStatusMessage struct {
	Type string // result, finish
	Data []byte
}

func (u *Flow) RunFlow(ctx context.Context, flowId int64, params map[string]interface{}, parallel int) (runId string, err error) {
	flow, exist, err := u.flowRepo.GetFlowById(ctx, flowId)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", fmt.Errorf("flow not exist")
	}

	return u.RunFlowByDetail(ctx, flow, params, parallel)
}

func (u *Flow) RunFlowSync(ctx context.Context, flowId int64, params map[string]interface{}, parallel int, outputNodeId string) (rsp writeflow.Map, err error) {
	flow, exist, err := u.flowRepo.GetFlowById(ctx, flowId)
	if err != nil {
		return writeflow.Map{}, err
	}
	if !exist {
		return writeflow.Map{}, fmt.Errorf("flow not exist")
	}

	if outputNodeId != "" {
		flow.Graph.OutputNodeId = outputNodeId
	}

	return u.RunFlowByDetailSync(ctx, flow, params, parallel)
}

func (u *Flow) RunFlowByDetail(ctx context.Context, flow *model.Flow, params map[string]interface{}, parallel int) (runId string, err error) {
	f, err := model.FlowFromModel(flow)
	if err != nil {
		return "", err
	}

	//log.Infof("flow: %+v", flow)
	//log.Infof("f: %+v", f)

	runId = fmt.Sprintf("flow.%s", uuid.New().String())

	status, err := u.wirteflow.ExecFlowAsync(ctx, f, params, parallel)
	if err != nil {
		return "", err
	}

	// get status async
	log.Infof("%s start", runId)
	start := time.Now()
	go func() {
		defer func() {
			log.Infof("%s end, spend: %s", runId, time.Now().Sub(start))
		}()

		for r := range status {
			bs, err := r.Json()
			if err != nil {
				log.Errorf("status to json err: %v", err)
				return
			}
			//log.Infof("%s %s", runId, bs)
			err = u.ws.Send(runId, bs)
			if err != nil {
				log.Errorf("ws send err: %v", err)
				return
			}
		}
		err = u.ws.Send(runId, ws.EOF)
		if err != nil {
			log.Errorf("ws send err: %v", err)
			return
		}
	}()

	return
}

func (u *Flow) RunFlowByDetailSync(ctx context.Context, flow *model.Flow, params map[string]interface{}, parallel int) (rsp writeflow.Map, err error) {
	f, err := model.FlowFromModel(flow)
	if err != nil {
		return writeflow.Map{}, err
	}

	return u.wirteflow.ExecNode(ctx, f, params, parallel)
}

func (u *Flow) AddWs(key string, conn *websocket.Conn) {
	u.ws.Add(key, conn)
}
