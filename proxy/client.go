package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/parinpan/protoevent"
)

type Proxy struct {
	ServerAddr string
	ServerPort int
	LocalPort  int
	Identifier string
	Forwader   net.Conn
}

func New(ServerAddr string, ServerPort int, LocalPort int, Identifier string) *Proxy {

	client := &Proxy{
		ServerAddr: ServerAddr,
		ServerPort: ServerPort,
		LocalPort:  LocalPort,
		Identifier: Identifier,
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

		go func() {
			defer client.Close()
			f := &Forwader{Id: p.Identifier, Conn: conn}
			_, err := io.Copy(f, p.Forwader)

			if err != nil {
				fmt.Println(err)
			}
		}()
	})

	event.OnMessageReceived(func(conn net.Conn, message []byte, rawMessage []byte) {
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
