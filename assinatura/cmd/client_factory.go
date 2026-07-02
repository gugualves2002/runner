package cmd

import "github.com/kyriosdata/runner/internal/cli"

var newCLIClient = func(port int) *cli.Client {
	return cli.NewClient(port)
}
