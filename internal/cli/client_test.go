package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Post_ReturnsError_WhenServerReturnsBadRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"payload inválido"}`))
	}))
	defer srv.Close()

	client := NewClientWithBaseURL(srv.URL)
	_, err := client.Post(context.Background(), "/sign", map[string]string{"data": ""})
	if err == nil {
		t.Fatal("esperava erro, mas recebeu nil")
	}
}

func TestClient_Post_ReturnsBodyOnSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("esperava POST, mas recebeu %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"signature":"abc"}`))
	}))
	defer srv.Close()

	client := NewClientWithBaseURL(srv.URL)
	body, err := client.Post(context.Background(), "/sign", map[string]string{"data": "Zm9v"})
	if err != nil {
		t.Fatalf("esperava sucesso, mas recebeu erro: %v", err)
	}
	if string(body) != `{"signature":"abc"}` {
		t.Fatalf("esperava body exato, mas recebeu %s", string(body))
	}
}

func TestClient_Post_ReturnsError_WhenServerReturnsPlainError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("erro inesperado"))
	}))
	defer srv.Close()

	client := NewClientWithBaseURL(srv.URL)
	_, err := client.Post(context.Background(), "/validate", map[string]string{"data": "Zm9v"})
	if err == nil {
		t.Fatal("esperava erro, mas recebeu nil")
	}
}

func TestClient_Post_ReturnsError_WhenPayloadCannotBeSerialized(t *testing.T) {
	client := NewClientWithBaseURL("http://example.invalid")
	_, err := client.Post(context.Background(), "/sign", func() {})
	if err == nil {
		t.Fatal("esperava erro de serialização, mas recebeu nil")
	}
}
