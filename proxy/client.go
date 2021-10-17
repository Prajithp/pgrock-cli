package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/parinpan/protoevent"
)

type Proxy struct {
	RequestChannel chan []string
	ServerAddr     string
	ServerPort     int
	LocalPort      int
	Identifier     string
	Forwader       net.Conn
	AppClient      *Forwader
}

func New(ReqCh chan []string, ServerAddr string, ServerPort int, LocalPort int, Identifier string) *Proxy {

	client := &Proxy{
		RequestChannel: ReqCh,
		ServerAddr:     ServerAddr,
		ServerPort:     ServerPort,
		LocalPort:      LocalPort,
		Identifier:     Identifier,
	}

	return client
}

func (p *Proxy) Run() {
	address := fmt.Sprintf("%s:%d", p.ServerAddr, p.ServerPort)
	agent, event := protoevent.CreateAgent("tcp", address)
	agent.SetDefaultReadSize(4096)

	event.OnConnectionAccepted(func(conn net.Conn) {
		clientAddr := fmt.Sprintf("127.0.0.1:%d", p.LocalPort)
		client, err := net.Dial("tcp", clientAddr)

		if err != nil {
			fmt.Printf("Dial failed: %v", err)
			response := NewErrorResponse(503, p.Identifier)
			conn.Write(response.toByteSlice())
			return
		}
		p.Forwader = client
		p.AppClient = &Forwader{Id: p.Identifier, Conn: conn, RequestChannel: p.RequestChannel}

		go func() {
			defer client.Close()
			_, err := io.Copy(p.AppClient, p.Forwader)

			if err != nil {
				fmt.Println(err)
			}
		}()
	})

	event.OnMessageReceived(func(conn net.Conn, message []byte, rawMessage []byte) {
		buf := bufio.NewReader(strings.NewReader(string(rawMessage)))
		req, _ := http.ReadRequest(buf)
		p.AppClient.Request = req

		_, err := p.Forwader.Write(message)
		if err != nil {
			fmt.Println("Failed to write message to local server")
			return
		}
	})

	_ = agent.Run(func(conn net.Conn) error {
		message := &ConnectionMeta{
			Type:  "acceptProxy",
			Bytes: p.Identifier,
		}
		response, _ := json.Marshal(message)
		_, err := conn.Write(response)
		return err
	})
}
