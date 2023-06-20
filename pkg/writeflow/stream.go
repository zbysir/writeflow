package writeflow

import (
	"io"
	"sync"
	"time"
)

// StreamResponse 是可以重复使用的 流，因为一个流可以被多个节点使用
type StreamResponse[T any] struct {
	lock  sync.Mutex
	datas []T
	err   error
	//c         chan T
	done      chan struct{}
	closeOnce sync.Once
}

func (s *StreamResponse[T]) Display() string {
	return ""
}

func NewSteamResponse[T any]() *StreamResponse[T] {
	return &StreamResponse[T]{done: make(chan struct{}, 1)}
}

type SteamResponseStr StreamResponse[string]

func (s *StreamResponse[T]) NewReader() Reader[T] {
	return &Read[T]{s: s}
}

type Read[T any] struct {
	s   *StreamResponse[T]
	idx int
}

func (r *Read[T]) Read() (T, error) {
	for {
		select {
		case <-r.s.done:
			if r.idx < len(r.s.datas) {
				t := r.s.datas[len(r.s.datas)-1]
				r.idx = len(r.s.datas)
				return t, nil
			} else {
				var t T
				return t, r.s.err
			}
		case <-time.After(time.Second / 20):
			// todo Use sync.Cond to Broadcast data change
			if r.idx < len(r.s.datas) {
				t := r.s.datas[len(r.s.datas)-1]
				r.idx = len(r.s.datas)
				return t, nil
			}
		}
	}
}

func (r *Read[T]) ReadAll() ([]T, error) {
	r.s.Wait()

	if r.s.err == io.EOF {
		return r.s.datas, nil
	}

	return r.s.datas, r.s.err
}

type Reader[T any] interface {
	Read() (T, error)
	ReadAll() ([]T, error)
}

func (s *StreamResponse[T]) Append(a T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	//s.c <- a
	s.datas = append(s.datas, a)
}

func (s *StreamResponse[T]) Close(e error) {
	s.closeOnce.Do(func() {
		s.err = e
		close(s.done)
	})
}

func (s *StreamResponse[T]) Wait() {
	select {
	case <-s.done:
	}
}
