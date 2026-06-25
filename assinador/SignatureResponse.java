package com.github.kyriosdata.assinador.model;

public class SignatureResponse {
    private String signature;

    public SignatureResponse(String signature) {
        this.signature = signature;
    }

    public String getSignature() { return signature; }
}