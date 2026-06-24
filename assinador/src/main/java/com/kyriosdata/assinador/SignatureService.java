package com.kyriosdata.assinador;

/**
 * Interface que define operações de assinatura digital.
 * Implementações podem simular ou executar operações reais.
 */
public interface SignatureService {

    /**
     * Simula a criação de uma assinatura digital.
     *
     * @param payload O conteúdo a ser assinado (pode ser um arquivo, documento, etc)
     * @param keyAlias O alias da chave privada a usar para assinatura
     * @return A assinatura simulada em formato Base64
     * @throws SignatureException Se houver erro na assinatura
     */
    String sign(String payload, String keyAlias) throws SignatureException;

    /**
     * Simula a validação de uma assinatura digital.
     *
     * @param payload O conteúdo que foi assinado
     * @param signature A assinatura em formato Base64
     * @param keyAlias O alias da chave pública a usar para validação
     * @return true se a assinatura é válida, false caso contrário
     * @throws SignatureException Se houver erro na validação
     */
    boolean validate(String payload, String signature, String keyAlias) throws SignatureException;
}
