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
	Id  int64  `json:"id"`
	Key string `json:"key"`

	Data NodeData `json:"data"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	//Edges []Edge `json:"edges"`
}

type Node struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Id       string `json:"id"`
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`
	Type             string   `json:"type"` // = component.Key
	Data             NodeData `json:"data"`
	PositionAbsolute struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"positionAbsolute"`
	Selected bool `json:"selected"`
	Dragging bool `json:"dragging"`
}

type NodeSource struct {
	Type     string `json:"type"`     // local / git
	CmdType  string `json:"cmd_type"` // go_script / go_pkg
	GitUrl   string `json:"gitUrl"`
	GoScript string `json:"go_script"`
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
	Label       string `json:"label"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	//BaseClasses  []string          `json:"baseClasses"`
	Inputs map[string]interface{} `json:"inputs"` // key -> response (node_id.output_key)
	Source NodeSource             `json:"source"`
	//FilePath     string            `json:"filePath"`
	InputAnchors  []NodeAnchor     `json:"inputAnchors"`
	InputParams   []NodeInputParam `json:"inputParams"`
	OutputAnchors []NodeAnchor     `json:"outputAnchors"`
	Selected      bool             `json:"selected"`
}
