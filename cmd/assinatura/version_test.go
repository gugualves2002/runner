package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	cmd.Dir = "."

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Erro ao executar comando: %v", err)
	}

	saida := strings.TrimSpace(out.String())
	if !strings.Contains(saida, "dev") {
		t.Errorf("Saída esperada contendo 'dev', recebido: '%s'", saida)
	}
}

func TestSignCommandWithoutPayload(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "sign", "--key-alias", "test")
	cmd.Dir = "."

	var out bytes.Buffer
	cmd.Stderr = &out

	// Deve falhar sem --payload
	if err := cmd.Run(); err == nil {
		t.Errorf("Esperava erro ao executar sign sem --payload")
	}
}

func TestSignCommandWithoutKeyAlias(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "sign", "--payload", "test")
	cmd.Dir = "."

	var out bytes.Buffer
	cmd.Stderr = &out

	// Deve falhar sem --key-alias
	if err := cmd.Run(); err == nil {
		t.Errorf("Esperava erro ao executar sign sem --key-alias")
	}
}

func TestHelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--help")
	cmd.Dir = "."

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Erro ao executar comando: %v", err)
	}

	saida := strings.TrimSpace(out.String())
	if !strings.Contains(saida, "CLI") && !strings.Contains(saida, "assinatura") {
		t.Errorf("Ajuda não contém informações esperadas: '%s'", saida)
	}
}
