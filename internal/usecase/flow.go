package usecase

import (
	"context"
	"fmt"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/pkg/modules/builtin"
	"github.com/zbysir/writeflow/pkg/writeflow"
)

type Flow struct {
	flowRepo repo.Flow
	//componentRepo repo.Component
	wirteflow *writeflow.WriteFlow
}

func NewFlow(flowRepo repo.Flow) *Flow {
	wirteflow := writeflow.NewWriteFlow(writeflow.WithModules(builtin.New()))

	return &Flow{
		flowRepo:  flowRepo,
		wirteflow: wirteflow,
	}
}

func (u *Flow) GetComponents(ctx context.Context) (cs []writeflow.CategoryWithComponent, err error) {
	cs = u.wirteflow.GetComponentList()
	return cs, nil
}

func (u *Flow) GetComponentByKey(ctx context.Context, key string) (cs model.Component, exist bool, err error) {
	cs, exist, err = u.wirteflow.GetComponentByKey(key)
	return
}

func (u *Flow) RunFlow(ctx context.Context, flowId int64, params map[string]interface{}) (r map[string]interface{}, err error) {
	flow, exist, err := u.flowRepo.GetFlowById(ctx, flowId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("flow not exist")
	}

	f, err := writeflow.FlowFromModel(flow)
	if err != nil {
		return nil, err
	}

	log.Infof("flow: %+v", flow)
	log.Infof("f: %+v", f)

	r, err = u.wirteflow.ExecFlow(ctx, f, params)
	if err != nil {
		return nil, err
	}

	return r, nil
}
