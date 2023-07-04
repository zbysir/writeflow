package ws

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zbysir/writeflow/internal/pkg/ttlpool"
	"io"
	"math/rand"
	"sync"
	"time"
)

// WsHub
// todo 优化：使用占位符逻辑，当 key 生成就放入，只有放入了 key 才能写消息，优化任意 topic 的消息都能写入造成内存泄漏。
type WsHub struct {
	conns   map[string]*websocket.Conn
	history *ttlpool.Pool[[]Message]
	l       sync.Mutex
}

type Message []byte

func NewHub() *WsHub {
	return &WsHub{
		conns:   map[string]*websocket.Conn{},
		l:       sync.Mutex{},
		history: ttlpool.NewPool[[]Message](),
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (h *WsHub) Add(key string, conn *websocket.Conn) {
	h.l.Lock()
	defer h.l.Unlock()
	if key == "" {
		key = fmt.Sprintf("%v", rand.Int63())
	}

	// close old conn
	if o, ok := h.conns[key]; ok {
		_ = o.Close()
	}

	h.conns[key] = conn

	// send history
	messages, _ := h.history.Get(key)
	for _, v := range messages {
		if bytes.Equal(v, EOF) {
			conn.Close()
			delete(h.conns, key)
			break
		}
		conn.WriteMessage(1, v)
	}
}

var EOF = []byte("EOF")

var historyTtl = time.Minute * 5

func (h *WsHub) Send(key string, body []byte) error {
	h.l.Lock()
	defer h.l.Unlock()
	h.history.Update(key, func(v []Message) ([]Message, time.Duration) {
		return append(v, body), historyTtl
	})

	// close conn
	if bytes.Equal(body, EOF) {
		if o, ok := h.conns[key]; ok {
			o.Close()
		}
		delete(h.conns, key)
		return nil
	}

	if o, ok := h.conns[key]; ok {
		err := o.WriteMessage(1, body)
		if err != nil {
			return err
		}
	}

	return nil
}

type sender func([]byte) error

func (s sender) Write(p []byte) (n int, err error) {
	return len(p), s(p)
}

func (s sender) Close() (err error) {
	return s(EOF)
}

func (h *WsHub) GetKeyWrite(key string) interface {
	io.Writer
	io.Closer
} {
	return sender(func(body []byte) error {
		return h.Send(key, body)
	})
}

func (h *WsHub) Close(key string) {
	h.l.Lock()
	defer h.l.Unlock()
	if o, ok := h.conns[key]; ok {
		o.Close()
	}
	delete(h.conns, key)
	return
}

func (h *WsHub) SendAll(body []byte) error {
	for _, o := range h.conns {
		err := o.WriteMessage(1, body)
		if err != nil {
			return err
		}
	}

	return nil
}
