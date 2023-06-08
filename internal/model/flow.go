package model

import (
	"strings"
	"time"
)

// Flow: 存储流程
// Component: 存储组件（定义组件 Schema）
// Node: 存储节点（组件实例）
// Cmder: Component 可以转换为 Cmder，用于执行。

type Flow struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Graph       Graph     `json:"graph"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Locales map[string]string

// 一般情况下一个 component 对应一个同名的 cmd。

type Component struct {
	Id       int64         `json:"id"`
	Key      string        `json:"key"`      // 组件类型，需要全局唯一
	Category string        `json:"category"` // category key
	Data     ComponentData `json:"data"`
}

type Graph struct {
	Nodes        Nodes  `json:"nodes"`
	OutputNodeId string `json:"output_node_id"` // 指定节点的输出为整个流程的输出，如果不指定，则将会找到 Nodes 中 Id = OUTPUT 的节点作为输出。
}

func (g *Graph) GetOutputNodeId() string {
	if g.OutputNodeId != "" {
		return g.OutputNodeId
	}
	n, ok := g.Nodes.FindById("OUTPUT")
	if ok {
		return n.Id
	}
	return "OUTPUT"
}

type Nodes []Node

// FindById 根据 id 找到 node
func (n Nodes) FindById(id string) (*Node, bool) {
	for _, v := range n {
		if v.Id == id {
			return &v, true
		}
	}
	return nil, false
}

type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Node 是 Component 的实例
type Node struct {
	Id               string       `json:"id"` // 前端自己生成，随意，保证画布中不重复就行。
	Width            int          `json:"width"`
	Height           int          `json:"height"`
	Position         NodePosition `json:"position"`
	Type             string       `json:"type"` // = Component.Key
	Data             NodeData     `json:"data"`
	PositionAbsolute NodePosition `json:"position_absolute,omitempty"`
}

type ComponentGoScript struct {
	Script   string `json:"script,omitempty"`
	InputKey string `json:"input_key"`
}
type ComponentGoPackage struct {
	GitUrl string `json:"git_url,omitempty"`
}

type ComponentCmdType = string

const (
	GoScriptCmd  ComponentCmdType = "go_script"
	GoPackageCmd ComponentCmdType = "go_package"
	BuiltInCmd   ComponentCmdType = "builtin"
	NothingCmd   ComponentCmdType = "nothing"
)

// ComponentSource 组件数据源，可以用来得到 Cmd
type ComponentSource struct {
	CmdType    ComponentCmdType   `json:"cmd_type"` // go_script / git / builtin
	BuiltinCmd string             `json:"builtin_cmd"`
	GoPackage  ComponentGoPackage `json:"go_package,omitempty"`
	GoScript   ComponentGoScript  `json:"go_script,omitempty"`
}

type NodeInputAnchor struct {
	Id              string            `json:"id"`
	Name            map[string]string `json:"name"`
	Key             string            `json:"key"`
	Type            string            `json:"type"`           // 数据模型，如 string / int / any
	List            bool              `json:"list,omitempty"` // 是否是数组
	Optional        bool              `json:"optional,omitempty"`
	Dynamic         bool              `json:"dynamic,omitempty"` // 是否是动态输入，是动态输入才能删除。
	TargetNodeId    string            `json:"target_node_id"`    // 关联的节点 id
	TargetOutputKey string            `json:"target_output_key"` // 关联的节点输出 key
}

type NodeOutputAnchor struct {
	Id              string            `json:"id"`
	Name            map[string]string `json:"name"`
	Key             string            `json:"key"`
	Type            string            `json:"type"`           // 数据模型，如 string / int / any
	List            bool              `json:"list,omitempty"` // 是否是数组
	Dynamic         bool              `json:"dynamic,omitempty"` // 是否是动态输入，是动态输入才能删除。
}

type NodeInputParam struct {
	Id          string            `json:"id"`
	Name        map[string]string `json:"name"`
	Key         string            `json:"key"`
	Type        string            `json:"type"`         // 数据模型，如 string / int / any
	DisplayType string            `json:"display_type"` // 显示类型，如 code / text / textarea / select / checkbox / radio
	Optional    bool              `json:"optional,omitempty"`
	Dynamic     bool              `json:"dynamic,omitempty"` // 是否是动态输入，是动态输入才能删除。
	Value       string            `json:"value"`             // 输入的值
}

type NodeData = ComponentData

type ComponentData struct {
	Name          Locales           `json:"name"`
	Icon          string            `json:"icon"`
	Description   Locales           `json:"description"`
	Source        ComponentSource   `json:"source"`
	DynamicInput  bool              `json:"dynamic_input"`            // 是否可以添加动态输入
	DynamicOutput bool              `json:"dynamic_output"`           // 输出是否和动态输入一样
	InputAnchors  []NodeInputAnchor `json:"input_anchors,omitempty"`  // 输入锚点定义
	InputParams   []NodeInputParam  `json:"input_params,omitempty"`   // 字面参数定义
	OutputAnchors []NodeOutputAnchor `json:"output_anchors,omitempty"` // 输出锚点定义
	Inputs        map[string]string `json:"inputs"`                   // key -> response (node_id.output_key)
}

func (d *ComponentData) GetInputValue(key string) string {
	for _, v := range d.InputParams {
		if v.Key == key {
			if v.Value != "" {
				return v.Value
			}
			return d.Inputs[key]
		}
	}

	return ""
}

func (d *ComponentData) GetInputAnchorValue(key string) (nodeId string, outputKey string) {
	for _, v := range d.InputAnchors {
		if v.Key == key {
			if v.TargetNodeId != "" {
				return v.TargetNodeId, v.TargetOutputKey
			}

			ss := strings.Split(d.Inputs[key], ".")
			var nodeId string
			var outputKey string
			if len(ss) > 1 {
				nodeId = ss[0]
				outputKey = ss[1]
			}

			return nodeId, outputKey
		}
	}

	return "", ""
}
