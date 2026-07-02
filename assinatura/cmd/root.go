package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para gerenciar e usar o assinador.jar",
	Long: `O assinatura CLI é uma ferramenta para iniciar, parar e interagir
	com o serviço de assinatura 'assinador.jar', seja no modo local ou servidor.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(validateCmd)
}