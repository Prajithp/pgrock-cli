package client

import (
	"flag"
	"fmt"
	"os"
)

const usage1 string = `Usage: %s [OPTIONS]
Options:
`
const usage2 string = `
Example:
		pgrock -r 192.168.22.1 -p 1080 -l 5001
`

type Options struct {
	Remote    string
	Port      int
	LocalPort int
}

func ParseArgs() *Options {

	usage := func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage2)
	}

	flag.Usage = usage

	remote := flag.String(
		"remote",
		"127.0.0.1",
		"Pgrock Server Address",
	)

	port := flag.Int(
		"port",
		1080,
		"Port number of Pgrock Server",
	)

	localPort := flag.Int(
		"local-port",
		8080,
		"Port number of the app server",
	)

	flag.Parse()
	opts := &Options{
		Remote:    *remote,
		Port:      *port,
		LocalPort: *localPort,
	}

	if flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}

	return opts
}
