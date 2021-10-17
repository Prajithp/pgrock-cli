package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
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
	Conn           net.Conn
	Id             string
	RequestChannel chan []string
	Request        *http.Request
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

	if f.Request != nil {
		respBuf := bufio.NewReader(strings.NewReader(string(m)))
		res, err := http.ReadResponse(respBuf, f.Request)
		if err == nil {
			defer res.Body.Close()

			statusCode := fmt.Sprintf("%v", res.StatusCode)
			method := f.Request.Method
			url := f.Request.RequestURI
			row := []string{method, statusCode, url}
			f.RequestChannel <- row
		}
	}

	return len(m), err
}
