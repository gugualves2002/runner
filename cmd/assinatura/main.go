package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Variável que será sobrescrita pelo ldflags no build do CI/CD
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para simulação de assinaturas digitais",
	Long:  "CLI multiplataforma para invocar assinador.jar e gerenciar operações de assinatura digital",
}

func main() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Criar uma assinatura digital",
		Long:  "Cria uma assinatura digital a partir de um payload usando uma chave privada",
		RunE:  handleSign,
	}

	var validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validar uma assinatura digital",
		Long:  "Valida uma assinatura digital usando uma chave pública",
		RunE:  handleValidate,
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Gerenciar servidor de assinatura",
		Long:  "Iniciar, parar ou verificar status do servidor HTTP de assinatura",
	}

	var serverStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Inicia o servidor de assinatura",
		Long:  "Inicia um servidor HTTP que expõe operações de assinatura",
		RunE:  handleServerStart,
	}

	var serverStopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Para o servidor de assinatura",
		Long:  "Encerra a instância do servidor HTTP de assinatura",
		RunE:  handleServerStop,
	}

	var serverStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Verifica status do servidor",
		Long:  "Verifica se o servidor de assinatura está em execução",
		RunE:  handleServerStatus,
	}

	// Flags para sign
	var payload, keyAlias, serverHost, serverPort string
	var useServer bool
	signCmd.Flags().StringVar(&payload, "payload", "", "Conteúdo a ser assinado (obrigatório)")
	signCmd.Flags().StringVar(&keyAlias, "key-alias", "", "Alias da chave privada (obrigatório)")
	signCmd.Flags().BoolVar(&useServer, "server", false, "Usar servidor HTTP em vez de modo local")
	signCmd.Flags().StringVar(&serverHost, "host", "localhost", "Host do servidor (padrão: localhost)")
	signCmd.Flags().StringVar(&serverPort, "port", "8080", "Porta do servidor (padrão: 8080)")
	signCmd.MarkFlagRequired("payload")
	signCmd.MarkFlagRequired("key-alias")

	// Flags para validate
	var validatePayload, signature, validateKeyAlias string
	var validateUseServer bool
	validateCmd.Flags().StringVar(&validatePayload, "payload", "", "Conteúdo original (obrigatório)")
	validateCmd.Flags().StringVar(&signature, "signature", "", "Assinatura em Base64 (obrigatório)")
	validateCmd.Flags().StringVar(&validateKeyAlias, "key-alias", "", "Alias da chave pública (obrigatório)")
	validateCmd.Flags().BoolVar(&validateUseServer, "server", false, "Usar servidor HTTP em vez de modo local")
	validateCmd.Flags().StringVar(&serverHost, "host", "localhost", "Host do servidor (padrão: localhost)")
	validateCmd.Flags().StringVar(&serverPort, "port", "8080", "Porta do servidor (padrão: 8080)")
	validateCmd.MarkFlagRequired("payload")
	validateCmd.MarkFlagRequired("signature")
	validateCmd.MarkFlagRequired("key-alias")

	// Flags para server start
	var startPort, startHost string
	serverStartCmd.Flags().StringVar(&startHost, "host", "localhost", "Host do servidor (padrão: localhost)")
	serverStartCmd.Flags().StringVar(&startPort, "port", "8080", "Porta do servidor (padrão: 8080)")

	// Flags para server status
	var statusHost, statusPort string
	serverStatusCmd.Flags().StringVar(&statusHost, "host", "localhost", "Host do servidor (padrão: localhost)")
	serverStatusCmd.Flags().StringVar(&statusPort, "port", "8080", "Porta do servidor (padrão: 8080)")

	// Flags para server stop
	var stopHost, stopPort string
	serverStopCmd.Flags().StringVar(&stopHost, "host", "localhost", "Host do servidor (padrão: localhost)")
	serverStopCmd.Flags().StringVar(&stopPort, "port", "8080", "Porta do servidor (padrão: 8080)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(serverStartCmd)
	serverCmd.AddCommand(serverStopCmd)
	serverCmd.AddCommand(serverStatusCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// handleSign executa a operação de assinatura
func handleSign(cmd *cobra.Command, args []string) error {
	payload, _ := cmd.Flags().GetString("payload")
	keyAlias, _ := cmd.Flags().GetString("key-alias")
	useServer, _ := cmd.Flags().GetBool("server")
	serverHost, _ := cmd.Flags().GetString("host")
	serverPort, _ := cmd.Flags().GetString("port")

	// Se usar servidor, tenta conectar
	if useServer || IsServerRunning(serverHost, serverPort) {
		response, err := InvokeServerSign(serverHost, serverPort, payload, keyAlias)
		if err == nil {
			output, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(output))
			return nil
		}
		if useServer {
			return err
		}
	}

	// Fallback para modo local
	jarPath := findAssinadorJar()
	if jarPath == "" {
		return fmt.Errorf("assinador.jar não encontrado")
	}

	javaCmd := exec.Command("java", "-jar", jarPath, "sign", "--payload", payload, "--key-alias", keyAlias)
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr

	return javaCmd.Run()
}

// handleValidate executa a operação de validação
func handleValidate(cmd *cobra.Command, args []string) error {
	payload, _ := cmd.Flags().GetString("payload")
	signature, _ := cmd.Flags().GetString("signature")
	keyAlias, _ := cmd.Flags().GetString("key-alias")
	useServer, _ := cmd.Flags().GetBool("server")
	serverHost, _ := cmd.Flags().GetString("host")
	serverPort, _ := cmd.Flags().GetString("port")

	// Se usar servidor, tenta conectar
	if useServer || IsServerRunning(serverHost, serverPort) {
		response, err := InvokeServerValidate(serverHost, serverPort, payload, signature, keyAlias)
		if err == nil {
			output, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(output))
			return nil
		}
		if useServer {
			return err
		}
	}

	// Fallback para modo local
	jarPath := findAssinadorJar()
	if jarPath == "" {
		return fmt.Errorf("assinador.jar não encontrado")
	}

	javaCmd := exec.Command("java", "-jar", jarPath, "validate", "--payload", payload, "--signature", signature, "--key-alias", keyAlias)
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr

	return javaCmd.Run()
}

// handleServerStart inicia o servidor HTTP
func handleServerStart(cmd *cobra.Command, args []string) error {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")

	jarPath := findAssinadorJar()
	if jarPath == "" {
		return fmt.Errorf("assinador.jar não encontrado")
	}

	javaCmd := exec.Command("java", "-jar", jarPath, "server", "--host", host, "--port", port)
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr

	fmt.Printf("Iniciando servidor em http://%s:%s\n", host, port)
	return javaCmd.Run()
}

// handleServerStop para o servidor HTTP
func handleServerStop(cmd *cobra.Command, args []string) error {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")

	if !IsServerRunning(host, port) {
		fmt.Println("Servidor não está em execução")
		return nil
	}

	jarPath := findAssinadorJar()
	if jarPath == "" {
		return fmt.Errorf("assinador.jar não encontrado")
	}

	javaCmd := exec.Command("java", "-jar", jarPath, "server", "stop", "--host", host, "--port", port)
	return javaCmd.Run()
}

// handleServerStatus verifica o status do servidor
func handleServerStatus(cmd *cobra.Command, args []string) error {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")

	if IsServerRunning(host, port) {
		fmt.Printf("Servidor em execução em http://%s:%s\n", host, port)
		return nil
	}

	fmt.Printf("Servidor não está em execução em http://%s:%s\n", host, port)
	return nil
}

// findAssinadorJar procura pelo assinador.jar em locais conhecidos
func findAssinadorJar() string {
	possiblePaths := []string{
		"assinador.jar",
		"./assinador/target/assinador.jar",
		"../assinador/target/assinador.jar",
		"../../assinador/target/assinador.jar",
	}

	// Também tenta procurar no ~/.hubsaude/
	homeDir, err := os.UserHomeDir()
	if err == nil {
		possiblePaths = append(possiblePaths, filepath.Join(homeDir, ".hubsaude", "assinador", "assinador.jar"))
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}

	return ""
}
