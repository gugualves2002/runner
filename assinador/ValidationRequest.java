package com.github.kyriosdata.assinador.model;

public class ValidationRequest {
    private String data;
    private String signature;
    private String algorithm;
    private String pkcs11ConfigPath;
    private String pin; // PIN is needed to load the keystore to get the certificate
    private String alias;

    public ValidationRequest() {}

    public String getData() { return data; }
    public String getSignature() { return signature; }
    public String getAlgorithm() { return algorithm; }
    public String getPkcs11ConfigPath() { return pkcs11ConfigPath; }
    public String getPin() { return pin; }
    public String getAlias() { return alias; }
}