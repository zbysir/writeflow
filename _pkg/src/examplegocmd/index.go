package examplegocmd

import (
	"context"
	"github.com/samber/lo"
	"strings"
)

func Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {
	x := lo.Map[interface{}, string](params, func(s interface{}, _ int) string { return s.(string) })
	return []interface{}{strings.Join(x, " + ")}, err

	return []interface{}{"1"}, nil
}
