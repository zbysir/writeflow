package writeflow

import (
	"encoding/json"
	"github.com/spf13/cast"
	"time"
)

type ComponentScript struct {
	Source string `json:"source"` // 如果设置则使用这个 source，否则使用 input 里的值
}

type ComponentGoPackage struct {
	GitUrl string `json:"git_url,omitempty"`
}

type ComponentCmdType string

const (
	GoScriptCmd   ComponentCmdType = "go_script"
	JavaScriptCmd ComponentCmdType = "js_script"
	GoPackageCmd  ComponentCmdType = "go_package"
	BuiltInCmd    ComponentCmdType = "builtin"
	NothingCmd    ComponentCmdType = "nothing"
)

// ComponentSource 组件数据源，可以用来得到 Cmd
type ComponentSource struct {
	CmdType    ComponentCmdType   `json:"cmd_type"` // go_script / git / builtin
	BuiltinCmd string             `json:"builtin_cmd"`
	GoPackage  ComponentGoPackage `json:"go_package,omitempty"`
	Script     ComponentScript    `json:"script,omitempty"`
}

type NodeOutputAnchor struct {
	Name    map[string]string `json:"name"`
	Key     string            `json:"key"`
	Type    string            `json:"type"`              // 数据模型，如 string / int / any
	List    bool              `json:"list,omitempty"`    // 是否是数组
	Dynamic bool              `json:"dynamic,omitempty"` // 是否是动态输入，是动态输入才能删除。
}

type NodeInputParam struct {
	Name        map[string]string  `json:"name"`
	Key         string             `json:"key"`
	InputType   NodeInputType      `json:"input_type"`
	Type        string             `json:"type"`               // 数据模型，如 string / int / json / any / bool
	DisplayType string             `json:"display_type"`       // 显示类型，如 code / input / textarea / select / checkbox / radio / password
	Options     []string           `json:"options"`            // 如果是 select / checkbox / radio，需要提供 options
	Optional    bool               `json:"optional,omitempty"` // 是否是可选的
	Dynamic     bool               `json:"dynamic,omitempty"`  // 是否是动态输入，是动态输入才能删除。
	Value       interface{}        `json:"value"`              // 输入的字面量
	List        bool               `json:"list"`               // 支持链接多个输入
	Anchors     []NodeAnchorTarget `json:"anchors"`
}

type ComponentData struct {
	Name          Locales            `json:"name"`
	Icon          string             `json:"icon"`
	Description   Locales            `json:"description"`
	Source        ComponentSource    `json:"source"`
	DynamicInput  bool               `json:"dynamic_input"`            // 是否可以添加动态输入
	DynamicOutput bool               `json:"dynamic_output"`           // 输出是否和动态输入一样
	InputParams   []NodeInputParam   `json:"input_params,omitempty"`   // 字面参数定义
	OutputAnchors []NodeOutputAnchor `json:"output_anchors,omitempty"` // 输出锚点定义
	Config        ComponentConfig    `json:"config"`
}

type ComponentConfig interface {
}

type PasswordString string

func (p PasswordString) Display() string {
	return "******"
}

func (d *ComponentData) GetInputValue(key string) interface{} {
	for _, v := range d.InputParams {
		if v.Key == key {
			if v.DisplayType == "password" {
				return PasswordString(cast.ToString(v.Value))
			}
			return v.Value
		}
	}

	return ""
}

func (d *ComponentData) GetInputAnchorValue(key string) ([]NodeAnchorTarget, bool) {
	for _, v := range d.InputParams {
		if v.Key == key && v.InputType == NodeInputAnchor {
			return v.Anchors, v.List
		}
	}

	return nil, false
}

type Locales map[string]string

type Category struct {
	Key  string  `json:"key"`
	Name Locales `json:"name"`
	Desc Locales `json:"desc"`
}

type NodeStatus = string

const (
	StatusRunning     NodeStatus = "running"
	StatusSuccess     NodeStatus = "success"
	StatusFailed      NodeStatus = "failed"
	StatusPending     NodeStatus = "pending"
	StatusUnreachable NodeStatus = "unreachable" // 被 if 分支忽略
)

// NodeStatusLog save node run result
type NodeStatusLog struct {
	NodeId string     `json:"node_id"`
	Status NodeStatus `json:"status"`
	// todo result has can't marshal type
	Error  string                 `json:"error,omitempty"`
	Result map[string]interface{} `json:"result,omitempty"`
	RunAt  time.Time              `json:"run_at"`
	EndAt  time.Time              `json:"end_at,omitempty"`
	Spend  string                 `json:"spend,omitempty"`
}

func NewNodeStatusLog(nodeId string, status NodeStatus, error string, result map[string]interface{}, runAt time.Time, endAt time.Time) NodeStatusLog {
	s := NodeStatusLog{
		NodeId: nodeId,
		Status: status,
		Error:  error,
		Result: result,
		RunAt:  runAt,
		EndAt:  endAt}

	if !s.EndAt.IsZero() {
		s.Spend = s.EndAt.Sub(s.RunAt).String()
	}
	return s
}

func (r *NodeStatusLog) Json() ([]byte, error) {
	// 过滤私密信息
	for k, v := range r.Result {
		if d, ok := v.(interface{ Display() string }); ok {
			r.Result[k] = d.Display()
		}
	}

	return json.Marshal(r)
}
