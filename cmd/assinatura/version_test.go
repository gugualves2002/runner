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
