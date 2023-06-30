package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
)

type System interface {
	GetSetting(ctx context.Context) (s *model.Setting, err error)
	SaveSetting(ctx context.Context, s *model.Setting) (err error)
}
