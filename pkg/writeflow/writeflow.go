package writeflow

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/modules"
	"github.com/zbysir/writeflow/pkg/schema"
)

type WriteFlow struct {
	option
	core *WriteFlowCore
}

func NewWriteFlow(ops ...Option) *WriteFlow {
	var o option
	for _, op := range ops {
		op(&o)
	}
	builtinCmd := map[string]schema.CMDer{}
	for _, m := range o.modules {
		//info := m.Info()
		for k, v := range m.Cmd() {
			//builtinCmd[info.NameSpace+"."+k] = v
			builtinCmd[k] = v
		}
	}
	c := NewWriteFlowCore()
	for _, v := range o.modules {
		for k, v := range v.Cmd() {
			c.RegisterCmd(k, v)
		}
	}

	return &WriteFlow{
		option: o,
		core:   c,
	}
}

type Option func(*option)

func WithModules(modules ...modules.Module) Option {
	return func(o *option) {
		o.modules = append(modules, o.modules...)
	}
}

type option struct {
	modules []modules.Module
}

type CategoryWithComponent struct {
	Category model.Category    `json:"category"`
	Children []model.Component `json:"children"`
}

func (w *WriteFlow) GetComponentList() []CategoryWithComponent {
	var components []model.Component
	var categories []model.Category
	for _, m := range w.option.modules {
		components = append(components, m.Components()...)
		categories = append(categories, m.Categories()...)
	}

	var componentByCategory = map[string][]model.Component{}
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

func (w *WriteFlow) GetComponentByKey(key string) (c model.Component, exist bool, err error) {
	var components []model.Component
	for _, m := range w.option.modules {
		components = append(components, m.Components()...)
	}

	for _, c := range components {
		if key == c.Type {
			return c, true, nil
		}
	}

	return c, false, nil
}

func (w *WriteFlow) ExecFlow(ctx context.Context, flow *Flow, initParams map[string]interface{}) (rsp map[string]interface{}, err error) {
	return w.core.ExecFlow(ctx, flow, initParams)
}

func (w *WriteFlow) ExecFlowAsync(ctx context.Context, flow *Flow, initParams map[string]interface{}) (status chan *model.NodeStatus, err error) {
	return w.core.ExecFlowAsync(ctx, flow, initParams)
}
