package client

import (
	"fmt"
	"net"

	"github.com/Prajithp/pgrock-cli/proxy"
	"github.com/parinpan/protoevent"
)

type Agent struct {
	ServerAddr string
	ServerPort int
	LocalPort  int
}

func NewAgent(ServerAddr string, ServerPort int, LocalPort int) *Agent {
	return &Agent{
		ServerAddr: ServerAddr,
		ServerPort: ServerPort,
		LocalPort:  LocalPort,
	}
}

func (a *Agent) Run() {
	serverAddress := fmt.Sprintf("%s:%d", a.ServerAddr, a.ServerPort)
	agent, event := protoevent.CreateAgent("tcp", serverAddress)
	agent.SetDefaultReadSize(4096)

	event.OnMessageReceived(func(conn net.Conn, message []byte, rawMessage []byte) {
		input, err := proxy.Parse(message)
		if err != nil {
			fmt.Println("Received a message: ", err)
			return
		}
		switch input.Type {
		case "init":
			fmt.Println("Connected to Server")
			fmt.Printf("Address %s", input.Bytes)
		case "reqProxy":
			proxy := proxy.New(a.ServerAddr, a.ServerPort, a.LocalPort, input.Bytes)
			proxy.Run()
		default:
			fmt.Printf("%s - %s", input.Type, input.Bytes)
		}
	})

	err := agent.Run(func(conn net.Conn) error {
		_, err := conn.Write([]byte(`{"type": "init", "bytes": "auth"}`))
		return err
	})

	if nil != err {
		panic(err)
	}
}
