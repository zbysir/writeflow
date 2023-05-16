package examplegocmd

import (
	"context"
	"github.com/samber/lo"
	"github.com/zbysir/writeflow/pkg/schema"
	"strings"
)

type Cmd struct {
	config map[string]interface{}
}

func (c *Cmd) Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {
	x := lo.Map[interface{}, string](params, func(s interface{}, _ int) string { return s.(string) })
	return []interface{}{strings.Join(x, " + ")}, err
}

func (c *Cmd) Schema(ctx context.Context) schema.CMDSchema {
	return schema.CMDSchema{}
}

func New(config map[string]interface{}) (schema.CMDer, error) {
	return &Cmd{config: config}, nil
}
