package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SignatureClient é a interface para se comunicar com o servidor assinador.
type SignatureClient struct {
	BaseURL string
}

// NewSignatureClient cria um novo cliente.
func NewSignatureClient(port int) *SignatureClient {
	return &SignatureClient{
		BaseURL: fmt.Sprintf("http://localhost:%d/api", port),
	}
}

// Post faz uma requisição POST para um endpoint com um corpo (payload).
func (c *SignatureClient) Post(endpoint string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	resp, err := http.Post(c.BaseURL+endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o servidor. O servidor está em execução? (use 'assinatura start'): %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do servidor: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Tenta decodificar a mensagem de erro do servidor
		var errorResponse struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(responseBody, &errorResponse) == nil && errorResponse.Message != "" {
			return nil, fmt.Errorf("servidor retornou erro (%s): %s", resp.Status, errorResponse.Message)
		}
		return nil, fmt.Errorf("servidor retornou erro (%s): %s", resp.Status, string(responseBody))
	}

	return responseBody, nil
}