package com.kyriosdata.assinador;

import java.nio.charset.StandardCharsets;
import java.util.Base64;

/**
 * Implementação simulada de assinatura digital.
 * Não utiliza criptografia real, apenas retorna valores pré-construídos para testes.
 */
public class FakeSignatureService implements SignatureService {

    /**
     * Simula a criação de uma assinatura digital.
     * Retorna uma assinatura pré-construída em Base64.
     *
     * @param payload O conteúdo a ser assinado
     * @param keyAlias O alias da chave privada
     * @return Uma assinatura simulada
     * @throws SignatureException Se payload ou keyAlias forem inválidos
     */
    @Override
    public String sign(String payload, String keyAlias) throws SignatureException {
        if (payload == null || payload.trim().isEmpty()) {
            throw new SignatureException("Payload não pode ser vazio");
        }
        if (keyAlias == null || keyAlias.trim().isEmpty()) {
            throw new SignatureException("Key alias não pode ser vazio");
        }

        // Simula uma assinatura combinando payload e keyAlias
        String signatureData = String.format("FAKE_SIGNATURE_%s_%s_%d",
                payload.hashCode(),
                keyAlias.hashCode(),
                System.currentTimeMillis());

        return Base64.getEncoder()
                .encodeToString(signatureData.getBytes(StandardCharsets.UTF_8));
    }

    /**
     * Simula a validação de uma assinatura digital.
     * Para fins de teste, retorna true se a assinatura contém "FAKE_SIGNATURE".
     *
     * @param payload O conteúdo que foi assinado
     * @param signature A assinatura em Base64
     * @param keyAlias O alias da chave pública
     * @return true se a assinatura parece válida, false caso contrário
     * @throws SignatureException Se parâmetros forem inválidos
     */
    @Override
    public boolean validate(String payload, String signature, String keyAlias) throws SignatureException {
        if (payload == null || payload.trim().isEmpty()) {
            throw new SignatureException("Payload não pode ser vazio");
        }
        if (signature == null || signature.trim().isEmpty()) {
            throw new SignatureException("Signature não pode ser vazia");
        }
        if (keyAlias == null || keyAlias.trim().isEmpty()) {
            throw new SignatureException("Key alias não pode ser vazio");
        }

        try {
            // Decodifica a assinatura de Base64
            String decodedSignature = new String(
                    Base64.getDecoder().decode(signature),
                    StandardCharsets.UTF_8);

            // Simula validação verificando se começa com "FAKE_SIGNATURE"
            return decodedSignature.startsWith("FAKE_SIGNATURE_");
        } catch (IllegalArgumentException e) {
            throw new SignatureException("Assinatura inválida (não é Base64): " + e.getMessage());
        }
    }
}
