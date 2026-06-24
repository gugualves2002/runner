package jdk

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Manager gerencia detecção e provisionamento do JDK
type Manager struct {
	requiredVersion string
	jdkPath         string
}

// NewManager cria um novo gerenciador de JDK
func NewManager(requiredVersion string) *Manager {
	return &Manager{
		requiredVersion: requiredVersion,
	}
}

// EnsureJDK verifica se JDK está disponível, caso contrário baixa
func (m *Manager) EnsureJDK() error {
	// Primeiro tenta encontrar java no PATH
	if m.isJavaAvailable() {
		return nil
	}

	// Tenta encontrar em ~/.hubsaude/jdk
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível obter diretório home: %w", err)
	}

	jdkDir := filepath.Join(homeDir, ".hubsaude", "jdk")
	m.jdkPath = filepath.Join(jdkDir, "bin", "java")
	if runtime.GOOS == "windows" {
		m.jdkPath = filepath.Join(jdkDir, "bin", "java.exe")
	}

	if m.isJavaAtPathAvailable(m.jdkPath) {
		return nil
	}

	// Precisa baixar JDK
	return m.downloadJDK(jdkDir)
}

// JavaPath retorna o caminho para o executável java
func (m *Manager) JavaPath() string {
	if m.jdkPath != "" {
		return m.jdkPath
	}
	return "java"
}

// isJavaAvailable verifica se java está no PATH
func (m *Manager) isJavaAvailable() bool {
	_, err := exec.LookPath("java")
	return err == nil
}

// isJavaAtPathAvailable verifica se java existe no caminho específico
func (m *Manager) isJavaAtPathAvailable(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// downloadJDK baixa o JDK do Adoptium (Eclipse Temurin)
func (m *Manager) downloadJDK(targetDir string) error {
	fmt.Println("JDK não encontrado. Baixando...")

	// Determina URL de download baseado no SO e arquitetura
	downloadURL := m.getDownloadURL()
	if downloadURL == "" {
		return fmt.Errorf("plataforma não suportada para download automático de JDK")
	}

	// Cria diretório se não existir
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", targetDir, err)
	}

	// Baixa o JDK
	fmt.Printf("Baixando de: %s\n", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("erro ao baixar JDK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar JDK: status %d", resp.StatusCode)
	}

	// Salva o arquivo compactado
	zipPath := filepath.Join(targetDir, "jdk.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao salvar JDK: %w", err)
	}

	fmt.Println("JDK baixado com sucesso")
	fmt.Println("⚠️  Descompactação manual necessária (função de descompactação não implementada)")
	fmt.Printf("Por favor, descompacte %s em %s\n", zipPath, targetDir)

	return nil
}

// getDownloadURL retorna a URL de download do JDK baseado na plataforma
func (m *Manager) getDownloadURL() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// URL base do Adoptium para Java 21
	baseURL := "https://api.adoptium.net/v3/binary/latest/21/ga"

	switch osName {
	case "windows":
		if arch == "amd64" {
			return baseURL + "/windows/x64/jre/hotspot/normal/eclipse"
		}
	case "linux":
		if arch == "amd64" {
			return baseURL + "/linux/x64/jre/hotspot/normal/eclipse"
		}
	case "darwin":
		if arch == "amd64" {
			return baseURL + "/mac/x64/jre/hotspot/normal/eclipse"
		}
	}

	return ""
}

// Version retorna a versão do Java instalado
func (m *Manager) Version() (string, error) {
	javaPath := m.JavaPath()
	cmd := exec.Command(javaPath, "-version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Parse da versão
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "", fmt.Errorf("não foi possível obter versão do Java")
}
