package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para gerenciar e usar o assinador.jar",
	Long: `O assinatura CLI é uma ferramenta para iniciar, parar e interagir
com o serviço de assinatura 'assinador.jar', seja no modo local ou servidor.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(validateCmd)
}