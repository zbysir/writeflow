package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
)

type Flow interface {
	GetFlowById(ctx context.Context, id string) (flow *model.Flow, exist bool, err error)
	GetComponentByKeys(ctx context.Context, keys []string) (components map[string]*model.Component, err error)
}
