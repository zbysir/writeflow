package examplegocmd

import (
	"context"
	"github.com/samber/lo"
	"github.com/zbysir/writeflow/pkg/schema"
	"strings"
)

type Config struct {
	Name string
}

func CreateConfig() *Config {
	return &Config{
		Name: "bysir",
	}
}

type Cmd struct {
	config *Config
}

func (c *Cmd) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	x := lo.Map[interface{}, string](params["_args"].([]interface{}), func(s interface{}, _ int) string { return s.(string) })
	return map[string]interface{}{"default": strings.Join(x, " + ")}, err
}

func (c *Cmd) Schema() schema.CMDSchema {
	return schema.CMDSchema{
		Inputs: []schema.CMDSchemaParams{
			{
				Key:         "_args",
				Type:        "[]any",
				DescLocales: map[string]string{},
			},
		},
		Outputs: []schema.CMDSchemaParams{
			{

				Key:         "default",
				Type:        "string",
				DescLocales: map[string]string{},
			},
		},
		Name:        "ExampleGoCmd",
		DescLocales: map[string]string{},
	}
}

func New(c *Config) (schema.CMDer, error) {
	return &Cmd{config: c}, nil
}
