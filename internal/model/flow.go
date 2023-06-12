package model

import (
	"github.com/samber/lo"
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

func (f *Flow) Upgrade() *Flow {
	return &Flow{
		Id:          f.Id,
		Name:        f.Name,
		Description: f.Description,
		Graph: Graph{
			Nodes: lo.Map(f.Graph.Nodes, func(k Node, i int) Node {
				params := append(k.Data.InputParams)
				for _, a := range k.Data.InputAnchors {
					a.InputType = "anchor"
					params = append(params, a)
				}

				for i, p := range params {
					if p.InputType == "anchor" {
						ss := strings.Split(k.Data.Inputs[p.Key], ".")
						if len(ss) == 2 {
							nodeId := ss[0]
							key := ss[1]
							params[i].Anchors = []NodeAnchorTarget{
								{
									NodeId:    nodeId,
									OutputKey: key,
								},
							}
						}
					} else {
						if v, ok := k.Data.Inputs[p.Key]; ok {
							params[i].Value = v
						}
					}
				}

				return Node{
					Id:       k.Id,
					Width:    k.Width,
					Height:   k.Height,
					Position: k.Position,
					Type:     k.Type,
					Data: NodeData{
						Name:          k.Data.Name,
						Icon:          k.Data.Icon,
						Description:   k.Data.Description,
						Source:        k.Data.Source,
						DynamicInput:  k.Data.DynamicInput,
						DynamicOutput: k.Data.DynamicOutput,
						CanDisable:    k.Data.CanDisable,
						InputAnchors:  nil,
						InputParams:   params,
						OutputAnchors: k.Data.OutputAnchors,
						Inputs:        nil,
						Config:        k.Data.Config,
					},
				}
			}),
			OutputNodeId: "",
		},
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

type Locales map[string]string

// 一般情况下一个 component 对应一个同名的 cmd。

type Component struct {
	Id       int64         `json:"id"`
	Type     string        `json:"type"`     // 组件类型，需要全局唯一
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
	Id       string       `json:"id"` // 前端自己生成，随意，保证画布中不重复就行。
	Width    int          `json:"width"`
	Height   int          `json:"height"`
	Position NodePosition `json:"position"`
	Type     string       `json:"type"` // = Component.Type
	Data     NodeData     `json:"data"`
}

type ComponentScript struct {
	InputKey string `json:"input_key"`
}

type ComponentGoPackage struct {
	GitUrl string `json:"git_url,omitempty"`
}

type ComponentCmdType = string

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

type NodeInputType = string

const (
	NodeInputTypeAnchor  NodeInputType = "anchor"
	NodeInputTypeLiteral NodeInputType = "literal"
)

type NodeInputParam struct {
	Name        map[string]string  `json:"name"`
	Key         string             `json:"key"`
	InputType   NodeInputType      `json:"input_type"`
	Type        string             `json:"type"`               // 数据模型，如 string / int / json / any
	DisplayType string             `json:"display_type"`       // 显示类型，如 code / input / textarea / select / checkbox / radio
	Options     []string           `json:"options"`            // 如果是 select / checkbox / radio，需要提供 options
	Optional    bool               `json:"optional,omitempty"` // 是否是可选的
	Dynamic     bool               `json:"dynamic,omitempty"`  // 是否是动态输入，是动态输入才能删除。
	Value       string             `json:"value"`              // 输入的字面量
	List        bool               `json:"list"`               // 支持链接多个输入
	Anchors     []NodeAnchorTarget `json:"anchors"`
}

type NodeAnchorTarget struct {
	NodeId    string `json:"node_id"`    // 关联的节点 id
	OutputKey string `json:"output_key"` // 关联的节点输出 key
}

type NodeData = ComponentData

type ForItemNode struct {
	NodeId    string `json:"node_id"`
	InputKey  string `json:"input_key"`
	OutputKey string `json:"output_key"` // outputKey 可不填，默认等于 inputKey
}

type ComponentData struct {
	Name          Locales         `json:"name"`
	Icon          string          `json:"icon"`
	Description   Locales         `json:"description"`
	Source        ComponentSource `json:"source"`
	DynamicInput  bool            `json:"dynamic_input"`  // 是否可以添加动态输入
	DynamicOutput bool            `json:"dynamic_output"` // 输出是否和动态输入一样
	CanDisable    bool            `json:"can_disable"`    // 是否可以禁用，如果禁用则不会执行，可以禁用的组件上会有一个开关，当关闭时需要填写一个 key 为 _enable 的输入。
	// InputAnchors 将要废弃
	InputAnchors  []NodeInputParam   `json:"input_anchors,omitempty"`  // 输入锚点定义
	InputParams   []NodeInputParam   `json:"input_params,omitempty"`   // 字面参数定义
	OutputAnchors []NodeOutputAnchor `json:"output_anchors,omitempty"` // 输出锚点定义
	// Inputs 将要废弃
	Inputs map[string]string `json:"inputs"` // key -> response (node_id.output_key)
	Config ComponentConfig   `json:"config"`
}

type ComponentConfig interface {
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

func (d *ComponentData) GetInputAnchorValue(key string) ([]NodeAnchorTarget, bool) {
	for _, v := range d.InputParams {
		if v.Key == key && v.InputType == NodeInputTypeAnchor {
			return v.Anchors, v.List
		}
	}

	return nil, false
}
