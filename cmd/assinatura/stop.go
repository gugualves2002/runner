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
	port, _ := cmd.Flags().GetInt("port")
	// US-01.8: Interromper execução do assinador.jar
	if err := stopServer(port); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao parar o servidor na porta %d: %v\n", port, err)
		os.Exit(1)
	}
	fmt.Printf("Servidor na porta %d parado com sucesso.\n", port)
}

func stopServer(p int) error {
	pid, err := getPID(p)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("nenhum servidor registrado na porta %d", p)
		}
		return fmt.Errorf("erro ao ler PID: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		// Processo não encontrado, talvez já tenha sido encerrado. Limpar o arquivo de PID.
		os.Remove(getPIDFilePath(p))
		return fmt.Errorf("processo com PID %d não encontrado: %w", pid, err)
	}

	if err := process.Kill(); err != nil {
		// Mesmo com erro, tenta limpar o arquivo de PID
		os.Remove(getPIDFilePath(p))
		return fmt.Errorf("falha ao encerrar processo com PID %d: %w", pid, err)
	}

	// Limpa o arquivo de PID após o sucesso
	if err := os.Remove(getPIDFilePath(p)); err != nil {
		return fmt.Errorf("processo encerrado, mas falha ao limpar arquivo de PID: %w", err)
	}

	return nil
}