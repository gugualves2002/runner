package com.github.kyriosdata.assinador;

import java.util.Base64;
import java.util.UUID;

/**
 * Implementação de simulação para o serviço de assinatura.
 * US-02.1
 */
public class FakeSignatureService implements SignatureService {

    @Override
    public String sign(String data, String algorithm) {
        // Validação básica de parâmetros
        if (data == null || data.isEmpty() || algorithm == null || algorithm.isEmpty()) {
            throw new IllegalArgumentException("Dados e algoritmo não podem ser nulos ou vazios.");
        }

        // Simula a criação de uma assinatura retornando uma resposta pré-construída
        String fakeSignature = "dados=" + data + "|alg=" + algorithm + "|uuid=" + UUID.randomUUID();
        return Base64.getEncoder().encodeToString(fakeSignature.getBytes());
    }

    @Override
    public boolean validate(String data, String signature, String algorithm) {
        // Simula a validação, retornando um resultado pré-determinado
        return signature != null && !signature.isEmpty() && signature.startsWith("ZGFkb3M9"); // "dados=" em Base64
    }
}