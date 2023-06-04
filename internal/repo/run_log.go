package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
)

type RunLog interface {
	GetRunLogById(ctx context.Context, id int64) (flow *model.RunLog, exist bool, err error)
	CreateRunLog(ctx context.Context, component *model.RunLog) (err error)
	UpdateRunLog(ctx context.Context, component *model.RunLog) (err error)
	DeleteRunLog(ctx context.Context, id int64) (err error)
	GetRunLogList(ctx context.Context, component GetRunLogListParams) (fs []model.RunLog, total int, err error)
}

type GetRunLogListParams struct {
	FlowId int `json:"flow_id" form:"flow_id"`
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}
