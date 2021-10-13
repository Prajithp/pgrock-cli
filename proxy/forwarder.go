package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
)

type ConnectionMeta struct {
	Type  string `json:"type"`
	Bytes string `json:"bytes"`
}

type ProxyMeta struct {
	ID       string `json:"id"`
	Response []byte `json:"response"`
}

type ProxyResponse struct {
	Type  string     `json:"type"`
	Bytes *ProxyMeta `json:"bytes"`
}

type Forwader struct {
	Conn net.Conn
	Id   string
}

func (f *Forwader) Write(m []byte) (int, error) {
	meta := &ProxyMeta{
		ID:       f.Id,
		Response: m,
	}
	response := &ProxyResponse{
		Type:  "proxy",
		Bytes: meta,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response)

	data := buf.Bytes()
	n, err := f.Conn.Write(data)
	if err != nil {
		return n, err
	}
	if n != len(data) {
		return n, io.ErrShortWrite
	}

	return len(m), err
}
