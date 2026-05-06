package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Variável que será sobrescrita pelo ldflags no build do CI/CD
var version = "dev"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "assinatura",
		Short: "CLI para simulação de assinaturas digitais",
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
