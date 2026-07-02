package com.github.kyriosdata.assinador.model;

public class SignatureRequest {
    private String data;
    private String algorithm;
    private String pkcs11ConfigPath;
    private String pin;
    private String alias;

    // Jackson precisa de um construtor no-args
    public SignatureRequest() {}

    public String getData() { return data; }
    public String getAlgorithm() { return algorithm; }
    public String getPkcs11ConfigPath() { return pkcs11ConfigPath; }
    public String getPin() { return pin; }
    public String getAlias() { return alias; }
}