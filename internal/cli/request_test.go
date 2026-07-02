package cli

import "testing"

func TestSignRequestValidate_ValidPayload(t *testing.T) {
	req := SignRequest{
		Data:      "Zm9v",
		Algorithm: "RSA",
	}

	if err := req.Validate(); err != nil {
		t.Fatalf("esperava payload válido, mas recebeu erro: %v", err)
	}
}

func TestSignRequestValidate_EmptyData(t *testing.T) {
	req := SignRequest{Algorithm: "RSA"}

	if err := req.Validate(); err == nil {
		t.Fatal("esperava erro para dados vazios, mas recebeu nil")
	}
}

func TestSignRequestValidate_MissingAliasWithPkcs11Config(t *testing.T) {
	req := SignRequest{
		Data:             "Zm9v",
		Algorithm:        "RSA",
		Pkcs11ConfigPath: "/tmp/config",
	}

	if err := req.Validate(); err == nil {
		t.Fatal("esperava erro quando pkcs11-config é informado sem alias")
	}
}

func TestValidateRequestValidate_ValidPayload(t *testing.T) {
	req := ValidateRequest{
		Data:      "Zm9v",
		Signature: "c2lnbmF0dXJl",
		Algorithm: "RSA",
	}

	if err := req.Validate(); err != nil {
		t.Fatalf("esperava payload válido, mas recebeu erro: %v", err)
	}
}

func TestValidateRequestValidate_MissingSignature(t *testing.T) {
	req := ValidateRequest{
		Data:      "Zm9v",
		Algorithm: "RSA",
	}

	if err := req.Validate(); err == nil {
		t.Fatal("esperava erro para assinatura vazia, mas recebeu nil")
	}
}

func TestValidateRequestValidate_MissingAliasWithPkcs11Config(t *testing.T) {
	req := ValidateRequest{
		Data:             "Zm9v",
		Signature:        "c2lnbmF0dXJl",
		Algorithm:        "RSA",
		Pkcs11ConfigPath: "/tmp/config",
	}

	if err := req.Validate(); err == nil {
		t.Fatal("esperava erro quando pkcs11-config é informado sem alias")
	}
}
