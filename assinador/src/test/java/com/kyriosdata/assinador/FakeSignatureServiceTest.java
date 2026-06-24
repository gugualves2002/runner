package com.kyriosdata.assinador;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Testes unitários para FakeSignatureService.
 */
class FakeSignatureServiceTest {

    private FakeSignatureService service;

    @BeforeEach
    void setup() {
        service = new FakeSignatureService();
    }

    // Testes de sign()

    @Test
    void testSignWithValidParameters() throws SignatureException {
        String payload = "test document content";
        String keyAlias = "test-key";

        String signature = service.sign(payload, keyAlias);

        assertNotNull(signature);
        assertFalse(signature.isEmpty());
        assertTrue(signature.contains("==")); // Típico de Base64
    }

    @Test
    void testSignThrowsExceptionWithNullPayload() {
        assertThrows(SignatureException.class, () -> {
            service.sign(null, "test-key");
        });
    }

    @Test
    void testSignThrowsExceptionWithEmptyPayload() {
        assertThrows(SignatureException.class, () -> {
            service.sign("", "test-key");
        });
    }

    @Test
    void testSignThrowsExceptionWithNullKeyAlias() {
        assertThrows(SignatureException.class, () -> {
            service.sign("test payload", null);
        });
    }

    @Test
    void testSignThrowsExceptionWithEmptyKeyAlias() {
        assertThrows(SignatureException.class, () -> {
            service.sign("test payload", "");
        });
    }

    @Test
    void testSignDifferentPayloadsDifferentSignatures() throws SignatureException {
        String keyAlias = "test-key";

        String signature1 = service.sign("payload 1", keyAlias);
        String signature2 = service.sign("payload 2", keyAlias);

        assertNotEquals(signature1, signature2);
    }

    // Testes de validate()

    @Test
    void testValidateWithValidSignature() throws SignatureException {
        String payload = "test document";
        String keyAlias = "test-key";

        String signature = service.sign(payload, keyAlias);
        boolean isValid = service.validate(payload, signature, keyAlias);

        assertTrue(isValid);
    }

    @Test
    void testValidateThrowsExceptionWithNullPayload() throws SignatureException {
        String signature = service.sign("test", "key");

        assertThrows(SignatureException.class, () -> {
            service.validate(null, signature, "key");
        });
    }

    @Test
    void testValidateThrowsExceptionWithNullSignature() {
        assertThrows(SignatureException.class, () -> {
            service.validate("test", null, "key");
        });
    }

    @Test
    void testValidateThrowsExceptionWithNullKeyAlias() throws SignatureException {
        String signature = service.sign("test", "key");

        assertThrows(SignatureException.class, () -> {
            service.validate("test", signature, null);
        });
    }

    @Test
    void testValidateReturnsFalseWithInvalidSignature() throws SignatureException {
        String invalidSignature = "INVALID_BASE64_SIGNATURE_NOT_REAL";

        assertThrows(SignatureException.class, () -> {
            service.validate("payload", invalidSignature, "key");
        });
    }

    @Test
    void testValidateReturnsFalseWithMalformedBase64() {
        String malformedBase64 = "not@valid@base64!!!";

        assertThrows(SignatureException.class, () -> {
            service.validate("payload", malformedBase64, "key");
        });
    }
}
