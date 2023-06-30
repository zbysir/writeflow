package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/libkv/store"
	"github.com/zbysir/writeflow/internal/model"
)

type BoltDBSystem struct {
	store store.Store
}

func NewBoltDBSystem(store store.Store) *BoltDBSystem {
	return &BoltDBSystem{store: store}
}

func (b *BoltDBSystem) GetSetting(ctx context.Context) (s *model.Setting, err error) {
	s = &model.Setting{}
	kv, err := b.store.Get(fmt.Sprintf("system/setting"))
	if err != nil {
		if err == store.ErrKeyNotFound {
			return s, nil
		}
		return nil, err
	}

	if kv == nil {
		return
	}

	err = json.Unmarshal(kv.Value, s)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal error: %w", err)
		return nil, err
	}

	return
}

func (b *BoltDBSystem) SaveSetting(ctx context.Context, s *model.Setting) (err error) {
	es, err := b.GetSetting(ctx)
	if err != nil {
		return err
	}

	ss := es.Merge(*s)
	bs, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	err = b.store.Put("system/setting", bs, nil)
	if err != nil {
		return err
	}

	return nil
}

var _ System = (*BoltDBSystem)(nil)
