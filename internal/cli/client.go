package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client envia requisições HTTP ao servidor assinador.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient cria um cliente HTTP com timeout padrão.
func NewClient(port int) *Client {
	return &Client{
		BaseURL:    fmt.Sprintf("http://localhost:%d/api/v1", port),
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// NewClientWithBaseURL cria um cliente HTTP para testes ou hosts personalizados.
func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Post envia uma requisição JSON e retorna o corpo da resposta.
func (c *Client) Post(ctx context.Context, endpoint string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o servidor. O servidor está em execução? (use 'assinatura start'): %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do servidor: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
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
