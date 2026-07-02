package cmd

import (
	"bytes"

	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kyriosdata/runner/internal/cli"
)

func TestRunSign_ReturnsError_WhenServerReturnsBadRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "payload inválido"})
	}))
	defer testServer.Close()

	newCLIClient = func(port int) *cli.Client {
		return cli.NewClientWithBaseURL(testServer.URL)
	}

	rootCmd.SetArgs([]string{"sign", "Zm9v", "RSA", "--port", "7070"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("esperava erro no comando sign, mas recebeu nil")
	}
}

func TestRunValidate_ReturnsError_WhenServerReturnsInvalidPayload(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "assinatura inválida"})
	}))
	defer testServer.Close()

	newCLIClient = func(port int) *cli.Client {
		return cli.NewClientWithBaseURL(testServer.URL)
	}

	rootCmd.SetArgs([]string{"validate", "Zm9v", "c2lnbmF0dXJl", "RSA", "--port", "7070"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("esperava erro no comando validate, mas recebeu nil")
	}
}

func TestRunSign_SuccessfulRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload cli.SignRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("erro ao decodificar request body: %v", err)
		}
		if payload.Data != "Zm9v" || payload.Algorithm != "RSA" {
			t.Fatalf("esperava payload com dados e algoritmo corretos, mas recebeu %+v", payload)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"signature": "abc"})
	}))
	defer testServer.Close()

	newCLIClient = func(port int) *cli.Client {
		return cli.NewClientWithBaseURL(testServer.URL)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"sign", "Zm9v", "RSA", "--port", "7070"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("esperava sucesso no comando sign, mas recebeu erro: %v", err)
	}

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	if !bytes.Contains(output, []byte("Assinatura criada com sucesso:")) {
		t.Fatalf("esperava mensagem de sucesso, mas recebeu %s", output)
	}
}

func TestRunValidate_SuccessfulRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload cli.ValidateRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("erro ao decodificar request body: %v", err)
		}
		if payload.Data != "Zm9v" || payload.Signature != "c2lnbmF0dXJl" || payload.Algorithm != "RSA" {
			t.Fatalf("esperava payload correto, mas recebeu %+v", payload)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"valid": true})
	}))
	defer testServer.Close()

	newCLIClient = func(port int) *cli.Client {
		return cli.NewClientWithBaseURL(testServer.URL)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"validate", "Zm9v", "c2lnbmF0dXJl", "RSA", "--port", "7070"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("esperava sucesso no comando validate, mas recebeu erro: %v", err)
	}

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	if !bytes.Contains(output, []byte("Resultado da validação: true")) {
		t.Fatalf("esperava mensagem de validação verdadeira, mas recebeu %s", output)
	}
}

func TestFlagErrorHandlesInvalidFlag(t *testing.T) {
	err := flagError(errors.New("flag inválida"), "port")
	if err == nil || err.Error() != "falha ao ler flag port: flag inválida" {
		t.Fatalf("esperava erro formatado, mas recebeu %v", err)
	}
}
