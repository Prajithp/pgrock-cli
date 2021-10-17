package client

import (
	"fmt"
	"net"

	"github.com/Prajithp/pgrock-cli/proxy"
	"github.com/parinpan/protoevent"
)

type Agent struct {
	Term       *TermV2
	ServerAddr string
	ServerPort int
	LocalPort  int
	Quit       chan int
}

func NewAgent(ServerAddr string, ServerPort int, LocalPort int) *Agent {
	return &Agent{
		ServerAddr: ServerAddr,
		ServerPort: ServerPort,
		LocalPort:  LocalPort,
		Quit:       make(chan int),
	}
}

func (a *Agent) Run() error {
	serverAddress := fmt.Sprintf("%s:%d", a.ServerAddr, a.ServerPort)
	agent, event := protoevent.CreateAgent("tcp", serverAddress)
	agent.SetDefaultReadSize(4096)

	event.OnConnectionAccepted(func(conn net.Conn) {
		go func() {
			for {
				select {
				case <-a.Quit:
					conn.Close()
					return
				}
			}
		}()
		return
	})

	event.OnMessageReceived(func(conn net.Conn, message []byte, rawMessage []byte) {
		input, err := proxy.Parse(message)
		if err != nil {
			fmt.Println("Received a message: ", err)
			return
		}
		switch input.Type {
		case "init":
			a.Term.SetTunnelStatus("Online", input.Bytes)
		case "reqProxy":
			proxy := proxy.New(a.Term.RequestChannel, a.ServerAddr, a.ServerPort, a.LocalPort, input.Bytes)
			proxy.Run()
		default:
			fmt.Printf("%s - %s", input.Type, input.Bytes)
		}
	})

	err := agent.Run(func(conn net.Conn) error {
		_, err := conn.Write([]byte(`{"type": "init", "bytes": "auth"}`))
		return err
	})

	return err
}
