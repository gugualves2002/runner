package invoker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPInvoker chama o assinador.jar em modo servidor HTTP.
type HTTPInvoker struct {
	client *http.Client
	url    string
}

// NewHTTPInvoker cria um invoker para servidor HTTP.
func NewHTTPInvoker(port int) *HTTPInvoker {
	return &HTTPInvoker{
		client: &http.Client{Timeout: 5 * time.Second},
		url:    fmt.Sprintf("http://localhost:%d/api/v1", port),
	}
}

// Sign solicita assinatura ao servidor HTTP.
func (h *HTTPInvoker) Sign(ctx context.Context, request interface{}) (map[string]interface{}, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url+"/sign", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição sign: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar servidor de assinatura: %w", err)
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}

// Validate solicita validação ao servidor HTTP.
func (h *HTTPInvoker) Validate(ctx context.Context, request interface{}) (map[string]interface{}, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url+"/validate", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição validate: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar servidor de validação: %w", err)
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}

func parseResponse(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct { Message string `json:"message"` }
		if json.Unmarshal(body, &errorResponse) == nil && errorResponse.Message != "" {
			return nil, fmt.Errorf("servidor retornou erro (%s): %s", resp.Status, errorResponse.Message)
		}
		return nil, fmt.Errorf("servidor retornou erro (%s): %s", resp.Status, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return response, nil
}
