package com.kyriosdata.assinador;

/**
 * Classe que encapsula a resposta de uma operação de assinatura.
 */
public class SignatureResponse {
    private String signature;
    private String status;
    private String message;

    public SignatureResponse(String signature, String status, String message) {
        this.signature = signature;
        this.status = status;
        this.message = message;
    }

    // Getters
    public String getSignature() {
        return signature;
    }

    public String getStatus() {
        return status;
    }

    public String getMessage() {
        return message;
    }

    // Setters
    public void setSignature(String signature) {
        this.signature = signature;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
