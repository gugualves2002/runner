package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <dados_base64> <assinatura_base64> <algoritmo>",
	Short: "Valida uma assinatura digital",
	Long:  `Envia dados e uma assinatura para o assinador.jar para validação.`,
	Args:  cobra.ExactArgs(3),
	Run:   runValidate,
}

func init() {
	validateCmd.Flags().IntP("port", "p", 7070, "Porta do servidor assinador")
	validateCmd.Flags().String("pkcs11-config", "", "Caminho para o arquivo de configuração PKCS#11")
	validateCmd.Flags().String("pin", "", "PIN do dispositivo PKCS#11 (necessário para obter a chave pública)")
	validateCmd.Flags().String("alias", "", "Alias do certificado no dispositivo PKCS#11")
}

func runValidate(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetInt("port")
	exitOnError(err)
	pkcs11Config, err := cmd.Flags().GetString("pkcs11-config")
	exitOnError(err)
	pin, err := cmd.Flags().GetString("pin")
	exitOnError(err)
	alias, err := cmd.Flags().GetString("alias")
	exitOnError(err)

	requestPayload := struct {
		Data             string `json:"data"`
		Signature        string `json:"signature"`
		Algorithm        string `json:"algorithm"`
		Pkcs11ConfigPath string `json:"pkcs11ConfigPath,omitempty"`
		Pin              string `json:"pin,omitempty"`
		Alias            string `json:"alias,omitempty"`
	}{
		Data:             args[0],
		Signature:        args[1],
		Algorithm:        args[2],
		Pkcs11ConfigPath: pkcs11Config,
		Pin:              pin,
		Alias:            alias,
	}

	client := NewSignatureClient(port)
	responseBody, err := client.Post("/validate", requestPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	var response struct {
		Valid bool `json:"valid"`
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao decodificar resposta do servidor: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Resultado da validação: %t\n", response.Valid)
}