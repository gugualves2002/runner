package cmd

import (
	"fmt"
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
	Run:   runStart,
}

var (
	port    int
	jarPath string
	timeout int
)

func init() {
	startCmd.Flags().IntVarP(&port, "port", "p", 7070, "Porta para o servidor escutar")
	startCmd.Flags().StringVar(&jarPath, "jar", "assinador.jar", "Caminho para o arquivo assinador.jar")
	startCmd.Flags().IntVar(&timeout, "timeout", 30, "Tempo em minutos para desligamento automático por inatividade (0 para desativar)")
}

func runStart(cmd *cobra.Command, args []string) {
	// US-01.7: Detectar instância em execução
	if isServerRunning(port) {
		fmt.Printf("Servidor já está em execução na porta %d.\n", port)
		return
	}

	// US-01.5: Iniciar assinador.jar no modo servidor
	fmt.Printf("Iniciando servidor na porta %d...\n", port)

	// Constrói os argumentos para o java -jar
	javaArgs := []string{
		"-jar",
		jarPath,
		"server",
		"--timeout",
		strconv.Itoa(timeout),
	}

	// O comando 'java' deve estar no PATH
	command := exec.Command("java", javaArgs...)

	// Redireciona a saída para arquivos de log para não poluir o terminal
	logDir := getConfigDir()
	os.MkdirAll(logDir, 0755)
	stdout, _ := os.Create(filepath.Join(logDir, fmt.Sprintf("assinador-%d.log", port)))
	stderr, _ := os.Create(filepath.Join(logDir, fmt.Sprintf("assinador-%d.err", port)))
	command.Stdout = stdout
	command.Stderr = stderr

	if err := command.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar o servidor: %v\n", err)
		fmt.Fprintf(os.Stderr, "Verifique se o Java está instalado e no PATH, e se o caminho para '%s' está correto.\n", jarPath)
		os.Exit(1)
	}

	pid := command.Process.Pid
	if err := savePID(port, pid); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar o PID: %v\n", err)
		// Tenta matar o processo que acabamos de iniciar
		if p, err := os.FindProcess(pid); err == nil {
			p.Kill()
		}
		os.Exit(1)
	}

	fmt.Printf("Servidor iniciado com sucesso! (PID: %d)\n", pid)
	fmt.Println("Aguardando o servidor ficar pronto...")

	// Aguarda o servidor responder ao health check
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if isServerRunning(port) {
			fmt.Println("Servidor pronto para receber requisições.")
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Fprintln(os.Stderr, "Erro: O servidor não respondeu a tempo.")
	stopServer(port) // Tenta limpar
	os.Exit(1)
}

func isServerRunning(p int) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/health", p))
	return err == nil && resp.StatusCode == http.StatusOK
}