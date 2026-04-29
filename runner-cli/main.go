package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type ReleaseInfo struct {
	Jar JarInfo `json:"jar"`
	Jre JreInfo `json:"jre"`
}

type JarInfo struct {
	URL     string `json:"url"`
	Version string `json:"version"`
}

type JreInfo struct {
	WindowsX64 string `json:"windows_x64"`
	LinuxX64   string `json:"linux_x64"`
	MacX64     string `json:"mac_x64"`
}

const (
	HubSaudeDir = ".hubsaude"
	ReleaseURL  = "https://raw.githubusercontent.com/gugualves2002/runner/main/release.json"
)

func main() {
	fmt.Println("Iniciando Sistema Runner...")

	// Criar estrutura de diretórios
	err := setupDirectories()
	if err != nil {
		fmt.Printf("Erro ao configurar diretórios: %v\n", err)
		os.Exit(1)
	}

	// Obter informações de release
	release, err := fetchReleaseInfo()
	if err != nil {
		fmt.Printf("Erro ao buscar release.json: %v\n", err)
		os.Exit(1)
	}

	// Determinar a URL do JRE baseada no SO
	jreURL := getJreUrlForOS(release.Jre)
	if jreURL == "" {
		fmt.Println("Sistema operacional ou arquitetura não suportados para download automático.")
		os.Exit(1)
	}

	fmt.Printf("Ambiente verificado. URL do JRE alvo: %s\n", jreURL)
}

// setupDirectories cria a pasta .hubsaude e subpastas na HOME do usuário
func setupDirectories() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	basePath := filepath.Join(homeDir, HubSaudeDir)
	dirs := []string{
		filepath.Join(basePath, "bin"),
		filepath.Join(basePath, "jre"),
		filepath.Join(basePath, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	fmt.Printf("Estrutura de diretórios configurada em: %s\n", basePath)
	return nil
}

func fetchReleaseInfo() (*ReleaseInfo, error) {

	// teste com mock, depois substituir pelo fetch real do url
	mockJSON := `{
		"jar": {
		  "url": "https://github.com/kyriosdata/assinador/releases/latest/download/assinador.jar",
		  "version": "1.2.0"
		},
		"jre": {
		  "windows_x64": "https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jre/hotspot/normal/eclipse",
		  "linux_x64":   "https://api.adoptium.net/v3/binary/latest/21/ga/linux/x64/jre/hotspot/normal/eclipse",
		  "mac_x64":     "https://api.adoptium.net/v3/binary/latest/21/ga/mac/x64/jre/hotspot/normal/eclipse"
		}
	}`

	var release ReleaseInfo
	err := json.Unmarshal([]byte(mockJSON), &release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func getJreUrlForOS(jre JreInfo) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Focando em x64 (amd64) conforme especificação do trabalho
	if arch != "amd64" {
		return ""
	}

	switch osName {
	case "windows":
		return jre.WindowsX64
	case "linux":
		return jre.LinuxX64
	case "darwin": // macOS
		return jre.MacX64
	default:
		return ""
	}
}
