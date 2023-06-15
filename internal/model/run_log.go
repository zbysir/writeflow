package model

import (
	"github.com/zbysir/writeflow/pkg/writeflow"
	"time"
)

type RunLog struct {
	Id       int64                     `json:"id"`
	FlowId   int64                     `json:"flow_id"`
	Status   writeflow.NodeStatus      `json:"status"`
	Result   []writeflow.NodeStatusLog `json:"result"` // save all node run result, update each node run result update.
	CreateAt time.Time                 `json:"create_at"`
}
