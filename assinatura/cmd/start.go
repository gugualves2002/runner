package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor",
	Long:  `Inicia o assinador.jar como um processo em background, escutando por requisições HTTP.`,
	RunE:  runStart,
}

func init() {
	startCmd.Flags().IntP("port", "p", 7070, "Porta para o servidor escutar")
	startCmd.Flags().String("jar", "assinador.jar", "Caminho para o arquivo assinador.jar")
	startCmd.Flags().Int("timeout", 30, "Tempo em minutos para desligamento automático por inatividade (0 para desativar)")
}

func runStart(cmd *cobra.Command, args []string) error {
	port, err := cmd.Flags().GetInt("port")
	if err := flagError(err, "port"); err != nil {
		return err
	}

	if isServerRunning(port) {
		slog.Info("servidor já em execução", "port", port)
		fmt.Printf("Servidor já está em execução na porta %d.\n", port)
		return nil
	}

	jarPath, err := cmd.Flags().GetString("jar")
		if err := flagError(err, "jar"); err != nil {
			return err
	}
	timeout, err := cmd.Flags().GetInt("timeout")
		if err := flagError(err, "timeout"); err != nil {
			return err
	}

	slog.Info("iniciando servidor", "port", port, "jar", jarPath, "timeoutMinutes", timeout)

	javaArgs := []string{
		"-jar",
		jarPath,
		"server",
		"--timeout",
		strconv.Itoa(timeout),
	}

	command := exec.Command("java", javaArgs...)

	logDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("não foi possível obter diretório de configuração: %w", err)
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("não foi possível criar diretório de log: %w", err)
	}

	stdout, err := os.Create(filepath.Join(logDir, fmt.Sprintf("assinador-%d.log", port)))
	if err != nil {
		slog.Warn("não foi possível criar arquivo de log stdout", "error", err)
	}
	stderr, err := os.Create(filepath.Join(logDir, fmt.Sprintf("assinador-%d.err", port)))
	if err != nil {
		slog.Warn("não foi possível criar arquivo de log stderr", "error", err)
	}
	command.Stdout = stdout
	command.Stderr = stderr

	if err := command.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar o servidor: %w. verifique se o Java está instalado e se '%s' existe", err, jarPath)
	}

	pid := command.Process.Pid
	if err := savePID(port, pid); err != nil {
		if p, killErr := os.FindProcess(pid); killErr == nil {
			p.Kill()
		}
		return fmt.Errorf("erro ao salvar o PID: %w", err)
	}

	fmt.Printf("Servidor iniciado com sucesso! (PID: %d)\n", pid)
	fmt.Println("Aguardando o servidor ficar pronto...")

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if isServerRunning(port) {
			fmt.Println("Servidor pronto para receber requisições.")
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err := stopServer(port); err != nil {
		slog.Warn("erro ao limpar servidor após falha no health check", "error", err)
	}
	return fmt.Errorf("o servidor não respondeu a tempo")
}

func isServerRunning(p int) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/api/health", p))
	if err == nil {
		defer resp.Body.Close()
	}
	return err == nil && resp.StatusCode == http.StatusOK
}
