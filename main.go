package main

import (
	"fmt"

	"github.com/Prajithp/pgrock-cli/client"
)

func main() {
	opts := client.ParseArgs()

	agent := client.NewAgent(opts.Remote, opts.Port, opts.LocalPort)
	term := client.NewTermV2(agent)
	agent.Term = term

	go func() {
		defer func() {
			term.App.Stop()
			fmt.Println("Closed the connection unexpectedly")
		}()
		err := agent.Run()
		fmt.Println(err)
	}()

	term.Draw()
}
