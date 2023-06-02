package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
)

type Component interface {
	CreateComponent(ctx context.Context, component *model.Component) (err error)
	GetComponentByKeys(ctx context.Context, keys []string) (components map[string]*model.Component, err error)
	DeleteComponent(ctx context.Context, key string) (err error)
	GetComponentList(ctx context.Context, params GetFlowListParams) (cs []model.Component, total int, err error)
}

type Flow interface {
	GetFlowById(ctx context.Context, id int64) (flow *model.Flow, exist bool, err error)
	CreateFlow(ctx context.Context, component *model.Flow) (err error)
	UpdateFlow(ctx context.Context, component *model.Flow) (err error)
	DeleteFlow(ctx context.Context, id int64) (err error)
	GetFlowList(ctx context.Context, component GetFlowListParams) (fs []model.Flow, total int, err error)
}

type GetFlowListParams struct {
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}
