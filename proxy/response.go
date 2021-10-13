package proxy

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
)

var statusReason = map[int]string{
	200: "OK",
	400: "Bad Request",
	403: "Forbidden",
	404: "Not Found",
	502: "Bad Gateway",
	503: "Service Unavailable",
}

type Response struct {
	httpVersion string
	status      int
	headers     map[string]string
	body        string
}

func NewResponse() Response {
	var res Response
	res.headers = make(map[string]string)
	res.status = 200
	res.addServerHeaders("1.1")

	return res
}

func NewErrorResponse(code int, identifier string) Response {
	res := NewResponse()
	res.status = code

	meta := &ProxyMeta{
		ID:       identifier,
		Response: []byte(statusReason[code]),
	}
	r := &ProxyResponse{
		Type:  "proxy",
		Bytes: meta,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(r)
	res.body = string(buf.Bytes())
	return res
}

func (res *Response) addServerHeaders(httpVersion string) {
	res.headers["Date"] = time.Now().String()
	res.headers["Content-Length"] = strconv.Itoa(len([]byte(res.body)))
	res.headers["Server"] = "Pgrock"
	res.httpVersion = httpVersion
}

func (res *Response) toByteSlice() []byte {
	var ret string
	ret += res.httpVersion + " " + strconv.Itoa(res.status) + " " + statusReason[res.status] + "\r\n"
	for k, v := range res.headers {
		ret += k + ": " + v + "\r\n"
	}
	ret += "\r\n"
	ret += res.body
	return []byte(ret)
}
