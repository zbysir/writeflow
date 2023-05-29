package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
)

type Flow interface {
	CreateComponent(ctx context.Context, component *model.Component) (err error)
	GetComponentByKeys(ctx context.Context, keys []string) (components map[string]*model.Component, err error)
	GetComponentList(ctx context.Context, params GetFlowListParams) (cs []model.Component, total int64, err error)

	GetFlowById(ctx context.Context, id int64) (flow *model.Flow, exist bool, err error)
	CreateFlow(ctx context.Context, component *model.Flow) (err error)
	GetFlowList(ctx context.Context, component GetFlowListParams) (fs []model.Flow, total int64, err error)
}

type GetFlowListParams struct {
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}
