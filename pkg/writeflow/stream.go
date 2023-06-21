package writeflow

import (
	"io"
	"sync"
	"time"
)

// StreamResponse 是可以重复使用的 流，因为一个流可以被多个节点使用
type StreamResponse[T any] struct {
	lock      sync.Mutex // lock for slice append
	data      []T
	err       error
	done      chan struct{}
	closeOnce sync.Once
}

// Display 返回空，不能被序列化
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
	if r.idx < len(r.s.data) {
		t := r.s.data[r.idx]
		r.idx++
		return t, nil
	}

	for {
		select {
		case <-r.s.done:
			if r.idx < len(r.s.data) {
				t := r.s.data[r.idx]
				r.idx++
				return t, nil
			} else {
				var t T
				if r.s.err != nil {
					return t, r.s.err
				}
				return t, io.EOF
			}
		case <-time.After(time.Second / 20):
			if r.idx < len(r.s.data) {
				t := r.s.data[r.idx]
				r.idx++
				return t, nil
			}
		}
	}
}

func (r *Read[T]) ReadAll() ([]T, error) {
	r.s.Wait()

	return r.s.data, r.s.err
}

type Reader[T any] interface {
	Read() (T, error)
	ReadAll() ([]T, error)
}

func (s *StreamResponse[T]) Append(a T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = append(s.data, a)
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
