package com.kyriosdata.assinador;

/**
 * Exceção lançada quando há erro em operações de assinatura.
 */
public class SignatureException extends Exception {

    public SignatureException(String message) {
        super(message);
    }

    public SignatureException(String message, Throwable cause) {
        super(message, cause);
    }
}
