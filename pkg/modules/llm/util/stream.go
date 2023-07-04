package util

import (
	"github.com/zbysir/writeflow/pkg/export"
	"io"
	"sync"
	"time"
)

// StreamResponse 是可以重复使用的 流，因为一个流可以被多个节点使用
type StreamResponse struct {
	lock      sync.Mutex // lock for slice append
	data      []string
	err       error
	done      chan struct{}
	closeOnce sync.Once
}

// Display 返回空，不能被序列化
func (s *StreamResponse) Display() string {
	return ""
}

func NewSteamResponse() *StreamResponse {
	return &StreamResponse{done: make(chan struct{}, 1)}
}

var _ export.Stream = (*StreamResponse)(nil)

func (s *StreamResponse) NewReader() export.Reader {
	return &Read{s: s}
}

type Read struct {
	s   *StreamResponse
	idx int
}

func (r *Read) Read() (string, error) {
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
				var t string
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

func (r *Read) ReadAll() ([]string, error) {
	r.s.Wait()

	return r.s.data, r.s.err
}

func (s *StreamResponse) Append(a string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = append(s.data, a)
}

func (s *StreamResponse) Close(e error) {
	s.closeOnce.Do(func() {
		s.err = e
		close(s.done)
	})
}

func (s *StreamResponse) Wait() {
	select {
	case <-s.done:
	}
}
