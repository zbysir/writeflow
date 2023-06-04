package model

import (
	"encoding/json"
	"time"
)

type Status = string

const (
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusPending Status = "pending"
)

// NodeStatus save node run result
type NodeStatus struct {
	NodeId string `json:"node_id"`
	Status Status `json:"status"`
	// todo result has can't marshal type
	Error  string                 `json:"error"`
	Result map[string]interface{} `json:"result"`
	RunAt  time.Time              `json:"run_at"`
	EndAt  time.Time              `json:"end_at"`
}

func (r *NodeStatus) Json() ([]byte, error) {
	return json.Marshal(r)
}
