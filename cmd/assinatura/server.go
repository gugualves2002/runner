package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultServerHost = "localhost"
	defaultServerPort = "8080"
)

// SignRequest representa a requisição para assinatura
type SignRequest struct {
	Payload  string `json:"payload"`
	KeyAlias string `json:"key-alias"`
}

// SignResponse representa a resposta de assinatura
type SignResponse struct {
	Signature string `json:"signature,omitempty"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// IsServerRunning verifica se o servidor está rodando
func IsServerRunning(host, port string) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("http://%s:%s/api/v1/health", host, port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// InvokeServerSign executa assinatura via servidor HTTP
func InvokeServerSign(host, port, payload, keyAlias string) (*SignResponse, error) {
	serverURL := fmt.Sprintf("http://%s:%s/api/v1/sign", host, port)

	data := url.Values{}
	data.Set("payload", payload)
	data.Set("key-alias", keyAlias)

	resp, err := http.PostForm(serverURL, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao servidor: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var response SignResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return &response, nil
}

// InvokeServerValidate executa validação via servidor HTTP
func InvokeServerValidate(host, port, payload, signature, keyAlias string) (*SignResponse, error) {
	serverURL := fmt.Sprintf("http://%s:%s/api/v1/validate", host, port)

	data := url.Values{}
	data.Set("payload", payload)
	data.Set("signature", signature)
	data.Set("key-alias", keyAlias)

	resp, err := http.PostForm(serverURL, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao servidor: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var response SignResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return &response, nil
}

// GetServerURL retorna a URL do servidor
func GetServerURL(host, port string) string {
	return fmt.Sprintf("http://%s:%s/api/v1", host, port)
}
