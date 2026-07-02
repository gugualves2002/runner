package com.github.kyriosdata.assinador;

public interface SignatureService {
    /**
     * Simula a criação de uma assinatura digital.
     */
    String sign(String data, String algorithm) throws Exception;

    boolean validate(String data, String signature, String algorithm) throws Exception;
}