package ttlpool

import (
	"sync"
	"time"
)

type TtlValue[T interface{}] struct {
	Value    T
	ExpireAt time.Time
}
type Pool[T interface{}] struct {
	l sync.Mutex

	values map[string]TtlValue[T]
}

func NewPool[T interface{}]() *Pool[T] {
	p := &Pool[T]{
		l:      sync.Mutex{},
		values: map[string]TtlValue[T]{},
	}
	p.cleanTtl()
	return p
}

func (p *Pool[T]) cleanTtl() {
	go func() {
		for {
			p.l.Lock()
			for k, v := range p.values {
				if v.ExpireAt.Before(time.Now()) {
					delete(p.values, k)
				}
			}
			p.l.Unlock()
			time.Sleep(time.Second)
		}
	}()
}

func (p *Pool[T]) Get(key string) (T, bool) {
	p.l.Lock()
	defer p.l.Unlock()

	v, ok := p.values[key]
	if !ok {
		var zero T
		return zero, false
	}
	if v.ExpireAt.Before(time.Now()) {
		delete(p.values, key)
		var zero T
		return zero, false
	}

	return v.Value, true
}
func (p *Pool[T]) Set(key string, value T, ttl time.Duration) {
	p.l.Lock()
	defer p.l.Unlock()

	p.values[key] = TtlValue[T]{
		Value:    value,
		ExpireAt: time.Now().Add(ttl),
	}
}

func (p *Pool[T]) Delete(key string) {
	p.l.Lock()
	defer p.l.Unlock()

	delete(p.values, key)
}

func (p *Pool[T]) Update(key string, fun func(t T) (T, time.Duration)) {
	p.l.Lock()
	defer p.l.Unlock()

	v, ok := p.values[key]
	if !ok || v.ExpireAt.Before(time.Now()) {
		var zero T
		t, duration := fun(zero)
		p.values[key] = TtlValue[T]{
			Value:    t,
			ExpireAt: time.Now().Add(duration),
		}
		return
	}

	t, duration := fun(v.Value)
	p.values[key] = TtlValue[T]{
		Value:    t,
		ExpireAt: time.Now().Add(duration),
	}
}
