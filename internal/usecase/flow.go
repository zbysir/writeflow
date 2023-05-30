package usecase

import (
	"context"
	"fmt"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/pkg/writeflow"
)

type Flow struct {
	flowRepo repo.Flow
}

func NewFlow(flowRepo repo.Flow) *Flow {
	return &Flow{flowRepo: flowRepo}
}

func (u *Flow) RunFlow(ctx context.Context, flowId int64, params map[string]interface{}) (r map[string]interface{}, err error) {
	flow, exist, err := u.flowRepo.GetFlowById(ctx, flowId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("flow not exist")
	}

	wf := writeflow.NewWriteFlow()

	f, err := writeflow.FlowFromModel(flow)
	if err != nil {
		return nil, err
	}

	log.Infof("flow: %+v", flow)
	log.Infof("f: %+v", f)
	componentKeys := f.UsedComponents()

	components, err := u.flowRepo.GetComponentByKeys(ctx, componentKeys)
	if err != nil {
		return nil, err
	}

	for _, c := range components {
		wc, err := writeflow.ComponentFromModel(c)
		if err != nil {
			return nil, err
		}

		wf.RegisterComponent(&wc)
	}

	r, err = wf.ExecFlow(ctx, f, params)
	if err != nil {
		return nil, err
	}

	return r, nil
}
