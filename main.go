package main

import (
	"fmt"
	"os"

	"github.com/Prajithp/pgrock-cli/client"
)

func main() {
	opts := client.ParseArgs()

	agent := client.NewAgent(opts.Remote, opts.Port, opts.LocalPort)

	term, err := client.NewScreen()
	if err != nil {
		fmt.Println("Could not initialize terminal")
		os.Exit(1)
	}
	term.Run(agent)
	agent.Run(term)
}
