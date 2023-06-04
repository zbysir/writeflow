package model

import (
	"time"
)

type RunLog struct {
	Id       int64        `json:"id"`
	FlowId   int64        `json:"flow_id"`
	Status   Status       `json:"status"`
	Result   []NodeStatus `json:"result"` // save all node run result, update each node run result update.
	CreateAt time.Time    `json:"create_at"`
}
