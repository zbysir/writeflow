package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/libkv/store"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/writeflow"
	"sort"
	"strconv"
	"time"
)

// BoltDBFlow 基于 boltdb 的实现，实际上它不是数据库，所以不支持太大的数据量。
type BoltDBFlow struct {
	store store.Store
}

func (b *BoltDBFlow) DeleteComponent(ctx context.Context, key string) (err error) {
	err = b.store.Delete(fmt.Sprintf("component/%v", key))
	if err != nil {
		return fmt.Errorf("store.Delete error: %w", err)
	}

	return nil
}

func (b *BoltDBFlow) DeleteFlow(ctx context.Context, id int64) (err error) {
	err = b.store.Delete(fmt.Sprintf("flow/%d", id))
	if err != nil {
		return fmt.Errorf("store.Delete error: %w", err)
	}

	return nil
}

func (b *BoltDBFlow) GetComponentList(ctx context.Context, params GetFlowListParams) (cs []writeflow.Component, total int, err error) {
	kv, err := b.store.List("component")
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("store.List error: %w", err)
	}

	for i, item := range kv {
		if i < params.Offset {
			continue
		}
		flow := writeflow.Component{}
		err = json.Unmarshal(item.Value, &flow)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal error: %w", err)
			return nil, 0, err
		}

		cs = append(cs, flow)
		if len(cs) >= params.Limit {
			break
		}
	}

	return cs, len(kv), nil
}

func (b *BoltDBFlow) GetFlowList(ctx context.Context, params GetFlowListParams) (fs []model.Flow, total int, err error) {
	kv, err := b.store.List("flow")
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	// kv 排序，按照 int 排序，大的在前
	sort.Slice(kv, func(i, j int) bool {
		ki, _ := strconv.ParseInt(kv[i].Key, 10, 64)
		kj, _ := strconv.ParseInt(kv[j].Key, 10, 64)
		if ki == kj {
			return kv[i].Key > kv[j].Key
		}

		return ki > kj
	})

	for i, item := range kv {
		if i < params.Offset {
			continue
		}
		flow := model.Flow{}
		err = json.Unmarshal(item.Value, &flow)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal error: %w", err)
			return nil, 0, err
		}
		if params.WithGraph {
			flow.Graph.Nodes = nil
		}

		fs = append(fs, flow)
		if len(fs) >= params.Limit {
			break
		}
	}

	return fs, len(kv), nil
}

func NewBoltDBFlow(store store.Store) *BoltDBFlow {
	return &BoltDBFlow{store: store}
}

var _ Flow = (*BoltDBFlow)(nil)

func (b *BoltDBFlow) IdSeq(namespace string) (id int64, err error) {
	// todo add lock
	kv, err := b.store.Get("id_seq/" + namespace)
	if err != nil {
		if err == store.ErrKeyNotFound {
			err = nil
		} else {
			return 0, fmt.Errorf("get id_seq error: %w", err)
		}
	}
	if kv != nil {
		err = json.Unmarshal(kv.Value, &id)
		if err != nil {
			return 0, err
		}
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
		if err == store.ErrKeyNotFound {
			return nil, false, nil
		}
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

	exist = true

	return
}

func (b *BoltDBFlow) GetComponentByKeys(ctx context.Context, keys []string) (components map[string]*model.Component, err error) {
	components = make(map[string]*model.Component)
	for _, key := range keys {
		kv, err := b.store.Get("component/key/" + key)
		if err != nil {
			if err == store.ErrKeyNotFound {
				err = nil
				continue
			}
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
	if component.Type == "" {
		return fmt.Errorf("key is empty")
	}

	bs, err := json.Marshal(component)
	if err != nil {
		return err
	}
	err = b.store.Put("component/key/"+component.Type, bs, nil)
	if err != nil {
		return err
	}

	return
}

func (b *BoltDBFlow) CreateFlow(ctx context.Context, fl *model.Flow) (id int64, err error) {
	id, err = b.IdSeq("flow")
	if err != nil {
		return 0, err
	}
	fl.Id = id
	now := time.Now()
	fl.UpdatedAt = now
	fl.CreatedAt = now
	bs, err := json.Marshal(fl)
	if err != nil {
		return 0, err
	}
	err = b.store.Put(fmt.Sprintf("flow/%v", id), bs, nil)
	if err != nil {
		return 0, err
	}

	return
}

func (b *BoltDBFlow) UpdateFlow(ctx context.Context, fl *model.Flow) (err error) {
	if fl.Id == 0 {
		return fmt.Errorf("id is empty")
	}
	existFlow, exist, err := b.GetFlowById(ctx, fl.Id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("flow not exist")
	}

	fl.CreatedAt = existFlow.CreatedAt

	bs, err := json.Marshal(fl)
	if err != nil {
		return err
	}
	err = b.store.Put(fmt.Sprintf("flow/%v", fl.Id), bs, nil)
	if err != nil {
		return err
	}

	return
}
