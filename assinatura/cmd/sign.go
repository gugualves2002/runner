package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign <dados_base64> <algoritmo>",
	Short: "Cria uma assinatura digital (simulada ou real)",
	Long:  `Envia dados para o assinador.jar para criar uma assinatura.`,
	Args:  cobra.ExactArgs(2),
	Run:   runSign,
}

func init() {
	signCmd.Flags().IntP("port", "p", 7070, "Porta do servidor assinador")
	signCmd.Flags().String("pkcs11-config", "", "Caminho para o arquivo de configuração PKCS#11")
	signCmd.Flags().String("pin", "", "PIN do dispositivo PKCS#11")
	signCmd.Flags().String("alias", "", "Alias da chave no dispositivo PKCS#11")
}

func runSign(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")
	pkcs11Config, _ := cmd.Flags().GetString("pkcs11-config")
	pin, _ := cmd.Flags().GetString("pin")
	alias, _ := cmd.Flags().GetString("alias")

	requestPayload := struct {
		Data             string `json:"data"`
		Algorithm        string `json:"algorithm"`
		Pkcs11ConfigPath string `json:"pkcs11ConfigPath,omitempty"`
		Pin              string `json:"pin,omitempty"`
		Alias            string `json:"alias,omitempty"`
	}{
		Data:             args[0],
		Algorithm:        args[1],
		Pkcs11ConfigPath: pkcs11Config,
		Pin:              pin,
		Alias:            alias,
	}

	client := NewSignatureClient(port)
	responseBody, err := client.Post("/sign", requestPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	var response struct {
		Signature string `json:"signature"`
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao decodificar resposta do servidor: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Assinatura criada com sucesso:")
	fmt.Println(response.Signature)
}