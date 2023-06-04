package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"math/rand"
	"sync"
	"time"
)

// WsHub todo 优化：使用占位符逻辑，当 key 生成就放入，只有放入了 key 才能写消息
type WsHub struct {
	conns   map[string]*websocket.Conn
	history map[string][]Message
	l       sync.Mutex
}

type Message []byte

func NewHub() *WsHub {
	return &WsHub{
		conns:   map[string]*websocket.Conn{},
		l:       sync.Mutex{},
		history: map[string][]Message{},
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
	if o, ok := h.conns[key]; ok {
		_ = o.Close()
	}

	h.conns[key] = conn

	// send history
	messages := h.history[key]
	for _, v := range messages {
		conn.WriteMessage(1, v)
	}
}

func (h *WsHub) Send(key string, body []byte) error {
	h.l.Lock()
	defer h.l.Unlock()

	h.history[key] = append(h.history[key], body)

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

func (h *WsHub) GetKeyWrite(key string) io.Writer {
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
	delete(h.history, key)
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
