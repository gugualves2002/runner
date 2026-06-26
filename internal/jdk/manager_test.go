package jdk

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager("21")
	if manager.requiredVersion != "21" {
		t.Errorf("Expected requiredVersion to be '21', but got '%s'", manager.requiredVersion)
	}
}

func TestEnsureJDK_Download(t *testing.T) {
	// 1. Prepara um servidor HTTP para simular o endpoint de download
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cria um arquivo zip falso para o teste
		zipBuffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(zipBuffer)
		// Adiciona um 'java' executável falso dentro do zip
		// A estrutura do zip do adoptium é algo como 'jdk-21.0.3+9-jre/bin/java'
		javaPathInZip := filepath.Join("jdk-21.0.3+9-jre", "bin", "java")
		if runtime.GOOS == "windows" {
			javaPathInZip += ".exe"
		}
		javaWriter, err := zipWriter.Create(javaPathInZip)
		if err != nil {
			t.Fatalf("Failed to create fake java in zip: %v", err)
		}
		javaWriter.Write([]byte("fake java executable")) // Conteúdo não importa
		zipWriter.Close()

		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipBuffer.Bytes())
	}))
	defer ts.Close()

	// 2. Prepara um ambiente de teste isolado
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome) // Para Windows

	// Isola o PATH para não encontrar o java do sistema
	t.Setenv("PATH", "")

	// 3. Executa a função
	manager := NewManager("21")
	// Sobrescreve a função que retorna a URL para apontar para nosso servidor de teste
	manager.getDownloadURLFunc = func() string { return ts.URL }

	err := manager.EnsureJDK()
	if err != nil {
		t.Fatalf("EnsureJDK() failed: %v", err)
	}

	// 4. Verifica o resultado
	// O caminho do JDK deve ter sido definido para o arquivo descompactado
	if manager.jdkPath == "" {
		t.Errorf("Expected jdkPath to be set, but it was empty")
	}

	// Verifica se o 'java' executável existe onde esperamos
	if _, err := os.Stat(manager.jdkPath); os.IsNotExist(err) {
		t.Errorf("Expected java executable at %s, but it does not exist", manager.jdkPath)
	}

	// Tenta executar o "java" falso (não é um executável real, então apenas verificamos o caminho)
	if !strings.Contains(manager.jdkPath, filepath.Join(".hubsaude", "jdk")) {
		t.Errorf("Expected jdkPath to be inside the .hubsaude/jdk directory, but got %s", manager.jdkPath)
	}
}
