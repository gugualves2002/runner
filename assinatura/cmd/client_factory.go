package cmd

import "github.com/gugualves2002/runner/internal/cli"

var newCLIClient = func(port int) *cli.Client {
	return cli.NewClient(port)
}
