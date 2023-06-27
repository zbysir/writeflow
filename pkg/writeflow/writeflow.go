package writeflow

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/zbysir/writeflow/pkg/plugin"
)

type WriteFlow struct {
	option
	modules []Module
	core    *WriteFlowCore
}

func (w *WriteFlow) RegisterModule(m Module) {
	w.modules = append(w.modules, m)
	for k, v := range m.Cmd() {
		w.core.RegisterCmd(k, v)
	}
}

type ModuleForPlugin struct {
	inner plugin.Plugin
}

func (m *ModuleForPlugin) Info() ModuleInfo {
	ii := m.inner.Info()
	return ModuleInfo{NameSpace: ii.NameSpace}
}

func (m *ModuleForPlugin) Categories() []Category {
	ic := m.inner.Categories()
	var cc []Category
	for _, c := range ic {
		cc = append(cc, Category{
			Key:  c.Key,
			Name: Locales(c.Name),
			Desc: Locales(c.Desc),
		})
	}

	return cc
}

func (m *ModuleForPlugin) Components() []Component {
	ic := m.inner.Components()
	var cc []Component
	for _, c := range ic {
		cc = append(cc, Component{
			Id:       c.Id,
			Type:     c.Type,
			Category: c.Category,
			Data: ComponentData{
				Name:        Locales(c.Data.Name),
				Icon:        c.Data.Icon,
				Description: Locales(c.Data.Description),
				Source: ComponentSource{
					CmdType:    ComponentCmdType(c.Data.Source.CmdType),
					BuiltinCmd: c.Data.Source.BuiltinCmd,
					Script:     ComponentScript{Source: c.Data.Source.Script},
				},
				DynamicInput:  c.Data.DynamicInput,
				DynamicOutput: c.Data.DynamicOutput,
				InputParams: lo.Map(c.Data.InputParams, func(item plugin.NodeInputParam, _ int) NodeInputParam {
					return NodeInputParam{
						Name:        Locales(item.Name),
						Key:         item.Key,
						InputType:   item.InputType,
						Type:        item.Type,
						DisplayType: item.DisplayType,
						Options:     item.Options,
						Optional:    item.Optional,
						Dynamic:     item.Dynamic,
						Value:       item.Value,
						List:        item.List,
						Anchors: lo.Map(item.Anchors, func(item plugin.NodeAnchorTarget, _ int) NodeAnchorTarget {
							return NodeAnchorTarget{
								NodeId:    item.NodeId,
								OutputKey: item.OutputKey,
							}
						}),
					}
				}),
				OutputAnchors: lo.Map(c.Data.OutputAnchors, func(item plugin.NodeOutputAnchor, _ int) NodeOutputAnchor {
					return NodeOutputAnchor{
						Name:    Locales(item.Name),
						Key:     item.Key,
						Type:    item.Type,
						Dynamic: item.Dynamic,
						List:    item.List,
					}
				}),
				Config: nil,
			},
		})
	}
	return cc
}

func (m *ModuleForPlugin) Cmd() map[string]CMDer {
	mm := map[string]CMDer{}
	for k, v := range m.inner.Cmd() {
		mm[k] = v
	}
	return mm
}

func (w *WriteFlow) RegisterPlugin(m plugin.Plugin) {
	w.modules = append(w.modules, &ModuleForPlugin{inner: m})
	for k, v := range m.Cmd() {
		w.core.RegisterCmd(k, v)
	}
}

func NewWriteFlow(ops ...Option) *WriteFlow {
	var o option
	for _, op := range ops {
		op(&o)
	}
	return &WriteFlow{
		option: o,
		core:   NewWriteFlowCore(),
	}
}

type Option func(*option)

type option struct {
}

type CategoryWithComponent struct {
	Category Category    `json:"category"`
	Children []Component `json:"children"`
}

type Component struct {
	Id       int64         `json:"id"`
	Type     string        `json:"type"`     // 组件类型，需要全局唯一
	Category string        `json:"category"` // category key
	Data     ComponentData `json:"data"`
}

func (w *WriteFlow) GetComponentList() []CategoryWithComponent {
	var components []Component
	var categories []Category
	for _, m := range w.modules {
		components = append(components, m.Components()...)
		categories = append(categories, m.Categories()...)
	}

	var componentByCategory = map[string][]Component{}
	for _, c := range components {
		componentByCategory[c.Category] = append(componentByCategory[c.Category], c)
	}

	var cwc []CategoryWithComponent
	for _, c := range categories {
		cwc = append(cwc, CategoryWithComponent{
			Category: c,
			Children: componentByCategory[c.Key],
		})
	}

	return cwc
}

func (w *WriteFlow) GetComponentByKey(key string) (c Component, exist bool, err error) {
	var components []Component
	for _, m := range w.modules {
		components = append(components, m.Components()...)
	}

	for _, c := range components {
		if key == c.Type {
			return c, true, nil
		}
	}

	return c, false, nil
}

func (w *WriteFlow) ExecNode(ctx context.Context, flow *Flow, initParams map[string]interface{}, parallel int) (rsp Map, err error) {
	return w.core.ExecNode(ctx, flow, initParams, parallel)
}

func (w *WriteFlow) ExecFlowAsync(ctx context.Context, flow *Flow, initParams map[string]interface{}, parallel int) (status chan NodeStatusLog, err error) {
	return w.core.ExecFlowAsync(ctx, flow, initParams, parallel)
}

type panicCmd struct {
	i CMDer
}

func (p *panicCmd) Exec(ctx context.Context, params Map) (rsp Map, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("cmd panic: %v", e)
		}
	}()

	return p.i.Exec(ctx, params)
}

func HandlePanicCmd(der CMDer) CMDer {
	return &panicCmd{i: der}
}
