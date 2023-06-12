package http_file_server

import (
	"bytes"
	"hash"
	"net/http"
)

type hashWriter struct {
	rw     http.ResponseWriter
	hash   hash.Hash
	buf    *bytes.Buffer
	len    int
	status int
}

func (hw hashWriter) Header() http.Header {
	return hw.rw.Header()
}

func (hw *hashWriter) WriteHeader(status int) {
	hw.status = status
}

func (hw *hashWriter) Write(b []byte) (int, error) {
	if hw.status == 0 {
		hw.status = http.StatusOK
	}
	// bytes.Buffer.Write(b) always return (len(b), nil), so just
	// ignore the return values.
	hw.buf.Write(b)

	l, err := hw.hash.Write(b)
	hw.len += l
	return l, err
}
