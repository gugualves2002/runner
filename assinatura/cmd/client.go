package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SignatureClient � a interface para se comunicar com o servidor assinador.
type SignatureClient struct {
	BaseURL string
	Client  http.Client
}

// NewSignatureClient cria um cliente HTTP para a porta padr�o do servidor.
func NewSignatureClient(port int) *SignatureClient {
	return &SignatureClient{
		BaseURL: fmt.Sprintf("http://localhost:%d/api/v1", port),
		Client:  http.Client{Timeout: 5 * time.Second},
	}
}

// Post faz uma requisi��o POST para um endpoint com um corpo JSON.
func (c *SignatureClient) Post(endpoint string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisio: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisicao HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o servidor. O servidor esta em execucao? (use 'assinatura start'): %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do servidor: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
	var errorResponse struct {
		Message string json:"message"
	}
	if json.Unmarshal(responseBody, &errorResponse) == nil && errorResponse.Message != "" {
		return nil, fmt.Errorf("servidor retornou erro (%s): %s", resp.Status, errorResponse.Message)
	}
		return nil, fmt.Errorf("servidor retornou erro (%s): %s", resp.Status, string(responseBody))
	}

	return responseBody, nil
}
