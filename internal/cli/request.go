package cli

import "fmt"

// SignRequest representa o payload de assinatura enviado ao servidor.
type SignRequest struct {
	Data             string `json:"data"`
	Algorithm        string `json:"algorithm"`
	Pkcs11ConfigPath string `json:"pkcs11ConfigPath,omitempty"`
	Pin              string `json:"pin,omitempty"`
	Alias            string `json:"alias,omitempty"`
}

// Validate valida os parâmetros de criação de assinatura.
func (r *SignRequest) Validate() error {
	if r.Data == "" {
		return fmt.Errorf("dados_base64 não pode ser vazio")
	}
	if r.Algorithm == "" {
		return fmt.Errorf("algoritmo não pode ser vazio")
	}
	if r.Alias == "" && r.Pkcs11ConfigPath != "" {
		return fmt.Errorf("alias é obrigatório quando pkcs11-config é fornecido")
	}
	return nil
}

// ValidateRequest representa o payload de validação enviado ao servidor.
type ValidateRequest struct {
	Data             string `json:"data"`
	Signature        string `json:"signature"`
	Algorithm        string `json:"algorithm"`
	Pkcs11ConfigPath string `json:"pkcs11ConfigPath,omitempty"`
	Pin              string `json:"pin,omitempty"`
	Alias            string `json:"alias,omitempty"`
}

// Validate valida os parâmetros de validação.
func (r *ValidateRequest) Validate() error {
	if r.Data == "" {
		return fmt.Errorf("dados_base64 não pode ser vazio")
	}
	if r.Signature == "" {
		return fmt.Errorf("assinatura_base64 não pode ser vazia")
	}
	if r.Algorithm == "" {
		return fmt.Errorf("algoritmo não pode ser vazio")
	}
	if r.Alias == "" && r.Pkcs11ConfigPath != "" {
		return fmt.Errorf("alias é obrigatório quando pkcs11-config é fornecido")
	}
	return nil
}
