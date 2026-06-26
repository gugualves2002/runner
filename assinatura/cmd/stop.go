package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para o assinador.jar em modo servidor",
	Long:  `Encerra o processo do assinador.jar que está em execução na porta especificada.`,
	Run:   runStop,
}

func init() {
	stopCmd.Flags().IntP("port", "p", 7070, "Porta do servidor a ser parado")
}

func runStop(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetInt("port")
	exitOnError(err)
	// US-01.8: Interromper execução do assinador.jar
	if err := stopServer(port); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao parar o servidor na porta %d: %v\n", port, err)
		os.Exit(1)
	}
	fmt.Printf("Servidor na porta %d parado com sucesso.\n", port)
}

func stopServer(p int) error {
	pid, err := getPID(p)
	// Tentamos obter o caminho do arquivo de PID para garantir a limpeza,
	// mesmo que a leitura do PID falhe.
	pidFile, pidFileErr := getPIDFilePath(p)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("nenhum servidor registrado na porta %d", p)
		}
		return fmt.Errorf("erro ao ler PID: %w", err)
	}

	// Sempre tentamos remover o arquivo de PID ao final da execução desta função.
	if pidFileErr == nil {
		defer os.Remove(pidFile)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("processo com PID %d não encontrado, arquivo de controle limpo: %w", pid, err)
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("falha ao encerrar processo com PID %d: %w", pid, err)
	}

	return nil
}