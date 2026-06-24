package invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/gugualves2002/runner/internal/jdk"
)

// Invoker gerencia a invocação do assinador.jar
type Invoker struct {
	jdkManager *jdk.Manager
	jarPath    string
	mode       string // "local" ou "http"
}

// NewInvoker cria um novo invoker
func NewInvoker(jarPath string, mode string) *Invoker {
	return &Invoker{
		jdkManager: jdk.NewManager("21"),
		jarPath:    jarPath,
		mode:       mode,
	}
}

// Sign executa operação de assinatura
func (i *Invoker) Sign(payload, keyAlias string) (map[string]interface{}, error) {
	// Garante que JDK está disponível
	if err := i.jdkManager.EnsureJDK(); err != nil {
		return nil, fmt.Errorf("erro ao preparar JDK: %w", err)
	}

	// Executa assinador.jar
	cmd := exec.Command(
		i.jdkManager.JavaPath(),
		"-jar", i.jarPath,
		"sign",
		"--payload", payload,
		"--key-alias", keyAlias,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro na assinatura: %s", stderr.String())
	}

	// Parse resposta JSON
	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return response, nil
}

// Validate executa operação de validação
func (i *Invoker) Validate(payload, signature, keyAlias string) (map[string]interface{}, error) {
	// Garante que JDK está disponível
	if err := i.jdkManager.EnsureJDK(); err != nil {
		return nil, fmt.Errorf("erro ao preparar JDK: %w", err)
	}

	// Executa assinador.jar
	cmd := exec.Command(
		i.jdkManager.JavaPath(),
		"-jar", i.jarPath,
		"validate",
		"--payload", payload,
		"--signature", signature,
		"--key-alias", keyAlias,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro na validação: %s", stderr.String())
	}

	// Parse resposta JSON
	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return response, nil
}

// FindJar procura pelo assinador.jar em locais conhecidos
func FindJar() string {
	possiblePaths := []string{
		"assinador.jar",
		"./assinador/target/assinador.jar",
		"../assinador/target/assinador.jar",
		"../../assinador/target/assinador.jar",
	}

	// Também tenta procurar no ~/.hubsaude/
	homeDir, err := os.UserHomeDir()
	if err == nil {
		possiblePaths = append(possiblePaths, fmt.Sprintf("%s/.hubsaude/assinador/assinador.jar", homeDir))
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}

	return ""
}

// FormatResponse formata a resposta para exibição no terminal
func FormatResponse(response map[string]interface{}, prettyPrint bool) string {
	if !prettyPrint {
		data, _ := json.Marshal(response)
		return string(data)
	}

	data, _ := json.MarshalIndent(response, "", "  ")
	return string(data)
}
