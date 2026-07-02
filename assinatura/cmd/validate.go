package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kyriosdata/runner/internal/cli"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <dados_base64> <assinatura_base64> <algoritmo>",
	Short: "Valida uma assinatura digital",
	Long:  `Envia dados e uma assinatura para o assinador.jar para validação.`,
	Args:  cobra.ExactArgs(3),
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().IntP("port", "p", 7070, "Porta do servidor assinador")
	validateCmd.Flags().String("pkcs11-config", "", "Caminho para o arquivo de configuração PKCS#11")
	validateCmd.Flags().String("pin", "", "PIN do dispositivo PKCS#11 (necessário para obter a chave pública)")
	validateCmd.Flags().String("alias", "", "Alias do certificado no dispositivo PKCS#11")
}

func runValidate(cmd *cobra.Command, args []string) error {
	port, err := cmd.Flags().GetInt("port")
	if err := flagError(err, "port"); err != nil {
		return err
	}
	pkcs11Config, err := cmd.Flags().GetString("pkcs11-config")
	if err := flagError(err, "pkcs11-config"); err != nil {
		return err
	}
	pin, err := cmd.Flags().GetString("pin")
	if err := flagError(err, "pin"); err != nil {
		return err
	}
	alias, err := cmd.Flags().GetString("alias")
	if err := flagError(err, "alias"); err != nil {
		return err
	}

	requestPayload := cli.ValidateRequest{
	Data:             args[0],
	Signature:        args[1],
	Algorithm:        args[2],
	Pkcs11ConfigPath: pkcs11Config,
	Pin:              pin,
	Alias:            alias,
	}
	if err := requestPayload.Validate(); err != nil {
		return err
	}

	slog.Info("executando comando validate", "port", port, "algorithm", requestPayload.Algorithm)

	client := newCLIClient(port)
	responseBody, err := client.Post(context.Background(), "/validate", requestPayload)
	if err != nil {
		return fmt.Errorf("falha ao validar assinatura: %w", err)
	}

	var response struct {
		Valid bool `json:"valid"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("erro ao decodificar resposta do servidor: %w", err)
	}

	fmt.Printf("Resultado da validação: %t\n", response.Valid)
	return nil
}
