package model

import (
	"github.com/zbysir/writeflow/pkg/writeflow"
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

type Component = writeflow.Component

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
	Id       string                  `json:"id"` // 前端自己生成，随意，保证画布中不重复就行。
	Width    int                     `json:"width"`
	Height   int                     `json:"height"`
	Position NodePosition            `json:"position"`
	Type     string                  `json:"type"` // = Component.Type
	Data     writeflow.ComponentData `json:"data"`
}

type NodeInputParam struct {
	Name        map[string]string            `json:"name"`
	Key         string                       `json:"key"`
	InputType   writeflow.NodeInputType      `json:"input_type"`
	Type        string                       `json:"type"`               // 数据模型，如 string / int / json / any
	DisplayType string                       `json:"display_type"`       // 显示类型，如 code / input / textarea / select / checkbox / radio / password
	Options     []string                     `json:"options"`            // 如果是 select / checkbox / radio，需要提供 options
	Optional    bool                         `json:"optional,omitempty"` // 是否是可选的
	Dynamic     bool                         `json:"dynamic,omitempty"`  // 是否是动态输入，是动态输入才能删除。
	Value       interface{}                  `json:"value"`              // 输入的字面量
	List        bool                         `json:"list"`               // 支持链接多个输入
	Anchors     []writeflow.NodeAnchorTarget `json:"anchors"`
}
