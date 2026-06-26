package cmd

import (
	"fmt"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro fatal ao processar comando: %v\n", err)
		os.Exit(1)
	}
}