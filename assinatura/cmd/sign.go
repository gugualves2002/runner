package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gugualves2002/runner/internal/cli"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign <dados_base64> <algoritmo>",
	Short: "Cria uma assinatura digital (simulada ou real)",
	Long:  `Envia dados para o assinador.jar para criar uma assinatura.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runSign,
}

func init() {
	signCmd.Flags().IntP("port", "p", 7070, "Porta do servidor assinador")
	signCmd.Flags().String("pkcs11-config", "", "Caminho para o arquivo de configuração PKCS#11")
	signCmd.Flags().String("pin", "", "PIN do dispositivo PKCS#11")
	signCmd.Flags().String("alias", "", "Alias da chave no dispositivo PKCS#11")
}

func runSign(cmd *cobra.Command, args []string) error {
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

	requestPayload := cli.SignRequest{
		Data:             args[0],
		Algorithm:        args[1],
		Pkcs11ConfigPath: pkcs11Config,
		Pin:              pin,
		Alias:            alias,
	}
	if err := requestPayload.Validate(); err != nil {
		return err
	}

	slog.Info("executando comando sign", "port", port, "algorithm", requestPayload.Algorithm)

	client := newCLIClient(port)
	responseBody, err := client.Post(context.Background(), "/sign", requestPayload)
	if err != nil {
		return fmt.Errorf("falha ao assinar: %w", err)
	}

	var response struct {
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("erro ao decodificar resposta do servidor: %w", err)
	}

	fmt.Println("Assinatura criada com sucesso:")
	fmt.Println(response.Signature)
	return nil
}
