package writeflow

import (
	"context"
	"fmt"
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
	c := NewWriteFlowCore()
	for _, v := range o.modules {
		for k, v := range v.Cmd() {
			//builtinCmd[info.NameSpace+"."+k] = v
			c.RegisterCmd(k, v)
		}
	}

	return &WriteFlow{
		option: o,
		core:   c,
	}
}

type Option func(*option)

func WithModules(modules ...Module) Option {
	return func(o *option) {
		o.modules = append(modules, o.modules...)
	}
}

type option struct {
	modules []Module
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
	for _, m := range w.option.modules {
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

func (w *WriteFlow) ExecFlow(ctx context.Context, flow *Flow, initParams map[string]interface{}, parallel int) (rsp map[string]interface{}, err error) {
	return w.core.ExecFlow(ctx, flow, initParams, parallel)
}

func (w *WriteFlow) ExecFlowAsync(ctx context.Context, flow *Flow, initParams map[string]interface{}, parallel int) (status chan NodeStatusLog, err error) {
	return w.core.ExecFlowAsync(ctx, flow, initParams, parallel)
}

type panicCmd struct {
	i schema.CMDer
}

func (p *panicCmd) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("cmd panic: %v", e)
		}
	}()

	return p.i.Exec(ctx, params)
}

func HandlePanicCmd(der schema.CMDer) schema.CMDer {
	return &panicCmd{i: der}
}
