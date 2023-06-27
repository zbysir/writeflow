// Package plugin should be kept simple to work well with yaegi
package export

import (
	"context"
)

type PluginInfo struct {
	NameSpace string
}

type Locales map[string]string

type Category struct {
	Key  string  `json:"key"`
	Name Locales `json:"name"`
	Desc Locales `json:"desc"`
}

type Component struct {
	Id       int64         `json:"id"`
	Type     string        `json:"type"`     // 组件类型，需要全局唯一
	Category string        `json:"category"` // category key
	Data     ComponentData `json:"data"`
}

type Plugin interface {
	Info() PluginInfo
	Categories() []Category
	Components() []Component
	Cmd() map[string]CMDer
}

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
}

type Register interface {
	RegisterPlugin(m Plugin)
}

// ComponentSource 组件数据源，可以用来得到 Cmd
type ComponentSource struct {
	CmdType    string `json:"cmd_type"` // go_script / git / builtin
	BuiltinCmd string `json:"builtin_cmd"`
	Script     string `json:"script,omitempty"`
}

type NodeOutputAnchor struct {
	Name    Locales `json:"name"`
	Key     string  `json:"key"`
	Type    string  `json:"type"`              // 数据模型，如 string / int / any
	List    bool    `json:"list,omitempty"`    // 是否是数组
	Dynamic bool    `json:"dynamic,omitempty"` // 是否是动态输入，是动态输入才能删除。
}

type NodeInputParam struct {
	Name        Locales            `json:"name"`
	Key         string             `json:"key"`
	InputType   string             `json:"input_type"`
	Type        string             `json:"type"`               // 数据模型，如 string / int / json / any / bool
	DisplayType string             `json:"display_type"`       // 显示类型，如 code / input / textarea / select / checkbox / radio / password
	Options     []string           `json:"options"`            // 如果是 select / checkbox / radio，需要提供 options
	Optional    bool               `json:"optional,omitempty"` // 是否是可选的
	Dynamic     bool               `json:"dynamic,omitempty"`  // 是否是动态输入，是动态输入才能删除。
	Value       interface{}        `json:"value"`              // 输入的字面量
	List        bool               `json:"list"`               // 支持链接多个输入
	Anchors     []NodeAnchorTarget `json:"anchors"`
}

type NodeAnchorTarget struct {
	NodeId    string `json:"node_id"`    // 关联的节点 id
	OutputKey string `json:"output_key"` // 关联的节点输出 key
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
}
