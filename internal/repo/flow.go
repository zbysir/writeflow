package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
)

type Flow interface {
	GetFlowById(ctx context.Context, id int64) (flow *model.Flow, exist bool, err error)
	GetComponentByKeys(ctx context.Context, keys []string) (components map[string]*model.Component, err error)
	CreateComponent(ctx context.Context, component *model.Component) (err error)
	CreateFlow(ctx context.Context, component *model.Flow) (err error)
}
