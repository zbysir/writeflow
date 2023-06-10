package writeflow

import (
	"context"
)

type nothingCMD struct {
}

func (n nothingCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return nil, nil
}

var _nothingCMD = nothingCMD{}
