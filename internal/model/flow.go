package model

import "time"

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
	X int `json:"x"`
	Y int `json:"y"`
}

// Node 是 Component 的实例
type Node struct {
	Id               string       `json:"id"` // 前段自己生成，随意，保证画布中不重复就行。
	Width            int          `json:"width"`
	Height           int          `json:"height"`
	Position         NodePosition `json:"position"`
	Type             string       `json:"type"` // = Component.Key
	Data             NodeData     `json:"data"`
	PositionAbsolute NodePosition `json:"positionAbsolute"`
}

type ComponentGoScript struct {
	Script string `json:"script"`
}

type ComponentSource struct {
	Type     string            `json:"type"`     // local / git
	CmdType  string            `json:"cmd_type"` // go_script / go_pkg
	GitUrl   string            `json:"gitUrl"`
	GoScript ComponentGoScript `json:"go_script"`
}

type NodeAnchor struct {
	Id   string            `json:"id"`
	Name map[string]string `json:"label"`
	Key  string            `json:"name"`
	Type string            `json:"type"`           // 数据模型，如 string / int / any
	List bool              `json:"list,omitempty"` // 是否是数组

}

type NodeInputParam struct {
	Id       string            `json:"id"`
	Name     map[string]string `json:"label"`
	Key      string            `json:"name"`
	Type     string            `json:"type"` // 数据模型，如 string / int / any
	Optional bool              `json:"optional,omitempty"`
}

type NodeData struct {
	ComponentData
	Inputs map[string]string `json:"inputs"` // key -> response (node_id.output_key)
}

type ComponentData struct {
	Name          map[string]string `json:"name"`
	Icon          string            `json:"icon"`
	Description   map[string]string `json:"description"`
	Source        ComponentSource   `json:"source"`
	InputAnchors  []NodeAnchor      `json:"inputAnchors"`  // 输入锚点定义
	InputParams   []NodeInputParam  `json:"inputParams"`   // 字面参数定义
	OutputAnchors []NodeAnchor      `json:"outputAnchors"` // 输出锚点定义
}
