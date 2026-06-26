package jdk

import (
	"archive/zip"
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
	requiredVersion    string
	jdkPath            string
	getDownloadURLFunc func() string
}

// NewManager cria um novo gerenciador de JDK
func NewManager(requiredVersion string) *Manager {
	m := &Manager{
		requiredVersion: requiredVersion,
	}
	m.getDownloadURLFunc = m.defaultGetDownloadURL
	return m
}

// EnsureJDK verifica se JDK está disponível, caso contrário baixa
func (m *Manager) EnsureJDK() error {
	// 1. Tenta encontrar 'java' no PATH do sistema
	if javaPath, err := exec.LookPath("java"); err == nil {
		if version, err := m.getJavaVersion(javaPath); err == nil && strings.HasPrefix(version, m.requiredVersion) {
			fmt.Println("JDK encontrado no PATH do sistema.")
			m.jdkPath = javaPath
			return nil
		}
	}

	// 2. Tenta encontrar JDK gerenciado localmente
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível obter diretório home: %w", err)
	}
	jdkDir := filepath.Join(homeDir, ".hubsaude", "jdk")

	if javaPath, err := m.findJavaExecutable(jdkDir); err == nil {
		if version, err := m.getJavaVersion(javaPath); err == nil && strings.HasPrefix(version, m.requiredVersion) {
			fmt.Println("JDK gerenciado encontrado localmente.")
			m.jdkPath = javaPath
			return nil
		}
	}

	// 3. Se não encontrou, baixa e instala
	fmt.Println("Nenhum JDK compatível encontrado. Baixando nova versão...")
	return m.downloadAndInstallJDK(jdkDir)
}

// JavaPath retorna o caminho para o executável java
func (m *Manager) JavaPath() string {
	if m.jdkPath != "" {
		return m.jdkPath
	}
	// Se nenhum JDK foi explicitamente definido (e.g. EnsureJDK não foi chamado),
	// assume que 'java' está no PATH.
	return "java"
}

func (m *Manager) getJavaVersion(javaPath string) (string, error) {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao executar 'java -version': %w\nOutput: %s", err, string(output))
	}

	versionOutput := string(output)
	lines := strings.Split(versionOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return strings.Trim(parts[2], `"`), nil
			}
		}
	}
	return "", fmt.Errorf("não foi possível parsear a versão do Java a partir de: %s", versionOutput)
}

func (m *Manager) findJavaExecutable(searchDir string) (string, error) {
	var javaPath string
	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// O nome do executável de Java
		javaExeName := "java"
		if runtime.GOOS == "windows" {
			javaExeName = "java.exe"
		}

		if !info.IsDir() && info.Name() == javaExeName {
			// Encontrou o executável, vamos parar de andar na árvore
			javaPath = path
			return io.EOF // Truque para parar o filepath.Walk
		}
		return nil
	})

	if err == io.EOF { // io.EOF é o erro que usamos para parar, não um erro real
		err = nil
	}

	if err != nil {
		return "", err
	}

	if javaPath == "" {
		return "", fmt.Errorf("executável 'java' não encontrado em %s", searchDir)
	}

	return javaPath, nil
}


func (m *Manager) downloadAndInstallJDK(targetDir string) error {
	downloadURL := m.getDownloadURL()
	if downloadURL == "" {
		return fmt.Errorf("plataforma não suportada para download automático de JDK (%s/%s)", runtime.GOOS, runtime.GOARCH)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", targetDir, err)
	}

	fmt.Printf("Baixando JDK de: %s\n", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("erro ao baixar JDK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar JDK (status: %s)", resp.Status)
	}

	// Usa um arquivo temporário para o download
	tmpFile, err := os.CreateTemp(targetDir, "jdk-*.zip")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Garante a remoção do zip no final

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao salvar JDK no arquivo temporário: %w", err)
	}
	tmpFile.Close() // Fecha para que o unzip possa abri-lo

	fmt.Println("Download completo. Descompactando...")
	if err := unzip(tmpFile.Name(), targetDir); err != nil {
		return fmt.Errorf("erro ao descompactar JDK: %w", err)
	}

	fmt.Println("JDK descompactado com sucesso.")

	// Após descompactar, encontra o caminho do java e o define no manager
	javaPath, err := m.findJavaExecutable(targetDir)
	if err != nil {
		return fmt.Errorf("JDK instalado, mas não foi possível encontrar o executável 'java': %w", err)
	}
	m.jdkPath = javaPath
	fmt.Printf("JDK configurado para usar: %s\n", m.jdkPath)

	return nil
}

func (m *Manager) getDownloadURL() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
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
		if arch == "amd64" { // Para Macs Apple Silicon (arm64), a URL seria outra
			return baseURL + "/mac/x64/jre/hotspot/normal/eclipse"
		}
	}
	return ""
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("falha ao abrir zip: %w", err)
	}
	defer r.Close()

	// O zip do JDK vem com um diretório raiz, e.g., "jdk-21.0.3+9-jre"
	// Vamos extrair o conteúdo preservando a estrutura.
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("caminho de arquivo inválido no zip: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("falha ao criar diretório para o arquivo: %w", err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("falha ao criar arquivo de destino: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("falha ao abrir arquivo dentro do zip: %w", err)
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("falha ao copiar conteúdo do arquivo: %w", err)
		}
	}
	return nil
}
