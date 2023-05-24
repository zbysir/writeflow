package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/libkv/store"
	"github.com/zbysir/writeflow/internal/model"
)

type BoltDBFlow struct {
	store store.Store
}

var _ Flow = (*BoltDBFlow)(nil)

func (b *BoltDBFlow) IdSeq(namespace string) (id int64, err error) {
	// todo add lock
	kv, err := b.store.Get("id_seq/" + namespace)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(kv.Value, &id)
	if err != nil {
		return 0, err
	}

	id = id + 1
	err = b.store.Put("id_seq/"+namespace, []byte(fmt.Sprintf("%v", id)), nil)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (b *BoltDBFlow) GetFlowById(ctx context.Context, id int64) (flow *model.Flow, exist bool, err error) {
	kv, err := b.store.Get(fmt.Sprintf("flow/%v", id))
	if err != nil {
		return nil, false, err
	}

	if kv == nil {
		return
	}

	flow = &model.Flow{}
	err = json.Unmarshal(kv.Value, flow)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal error: %w", err)
		return nil, false, err
	}

	return
}

func (b *BoltDBFlow) GetComponentByKeys(ctx context.Context, keys []string) (components map[string]*model.Component, err error) {
	components = make(map[string]*model.Component)
	for _, key := range keys {
		kv, err := b.store.Get("component/key/" + key)
		if err != nil {
			return nil, err
		}
		if kv == nil {
			continue
		}

		comp := &model.Component{}
		err = json.Unmarshal(kv.Value, comp)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal error: %w", err)
			return nil, err
		}

		components[key] = comp
	}

	return
}

func (b *BoltDBFlow) CreateComponent(ctx context.Context, component *model.Component) (err error) {
	if component.Key == "" {
		return fmt.Errorf("key is empty")
	}

	bs, err := json.Marshal(component)
	if err != nil {
		return err
	}
	err = b.store.Put("component/key/"+component.Key, bs, nil)
	if err != nil {
		return err
	}

	return
}

func (b *BoltDBFlow) CreateFlow(ctx context.Context, component *model.Flow) (err error) {
	bs, err := json.Marshal(component)
	if err != nil {
		return err
	}
	id, err := b.IdSeq("flow")
	if err != nil {
		return err
	}
	err = b.store.Put(fmt.Sprintf("flow/%v", id), bs, nil)
	if err != nil {
		return err
	}

	return
}
