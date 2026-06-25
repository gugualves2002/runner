package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// getConfigDir retorna o diretório de configuração (~/.runner)
func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Não foi possível obter o diretório home do usuário: " + err.Error())
	}
	return filepath.Join(home, ".runner")
}

// getPIDFilePath retorna o caminho para o arquivo de PID para uma porta específica
func getPIDFilePath(port int) string {
	return filepath.Join(getConfigDir(), fmt.Sprintf("assinador-%d.pid", port))
}

// savePID salva o PID em um arquivo
func savePID(port, pid int) error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(getPIDFilePath(port), []byte(strconv.Itoa(pid)), 0644)
}

// getPID lê o PID de um arquivo
func getPID(port int) (int, error) {
	data, err := os.ReadFile(getPIDFilePath(port))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}