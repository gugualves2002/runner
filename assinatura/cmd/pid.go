package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// getConfigDir retorna o diretório de configuração (~/.runner)
func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("não foi possível obter o diretório home do usuário: %w", err)
	}
	return filepath.Join(home, ".runner"), nil
}

// getPIDFilePath retorna o caminho para o arquivo de PID para uma porta específica
func getPIDFilePath(port int) (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, fmt.Sprintf("assinador-%d.pid", port)), nil
}

// savePID salva o PID em um arquivo
func savePID(port, pid int) error {
	pidFile, err := getPIDFilePath(port)
	if err != nil {
		return err
	}
	configDir := filepath.Dir(pidFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

// getPID lê o PID de um arquivo
func getPID(port int) (int, error) {
	pidFile, err := getPIDFilePath(port)
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}