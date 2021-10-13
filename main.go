package main

import (
	"github.com/Prajithp/pgrock/client"
)

func main() {
	opts := client.ParseArgs()

	agent := client.NewAgent(opts.Remote, opts.Port, opts.LocalPort)
	agent.Run()
}
