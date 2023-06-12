package model

import (
	"encoding/json"
	"time"
)

type Status = string

const (
	StatusRunning     Status = "running"
	StatusSuccess     Status = "success"
	StatusFailed      Status = "failed"
	StatusPending     Status = "pending"
	StatusUnreachable Status = "unreachable" // 被 if 分支忽略
)

// NodeStatus save node run result
type NodeStatus struct {
	NodeId string `json:"node_id"`
	Status Status `json:"status"`
	// todo result has can't marshal type
	Error  string                 `json:"error,omitempty"`
	Result map[string]interface{} `json:"result,omitempty"`
	RunAt  time.Time              `json:"run_at"`
	EndAt  time.Time              `json:"end_at,omitempty"`
	Spend  string                 `json:"spend,omitempty"`
}

func NewNodeStatus(nodeId string, status Status, error string, result map[string]interface{}, runAt time.Time, endAt time.Time) NodeStatus {
	s := NodeStatus{
		NodeId: nodeId,
		Status: status,
		Error:  error,
		Result: result,
		RunAt:  runAt,
		EndAt:  endAt}

	if s.EndAt.IsZero() {
		s.Spend = s.EndAt.Sub(s.RunAt).String()
	}
	return s
}

func (r *NodeStatus) Json() ([]byte, error) {
	return json.Marshal(r)
}
