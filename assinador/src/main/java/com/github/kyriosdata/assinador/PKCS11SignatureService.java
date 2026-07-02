package com.github.kyriosdata.assinador;

import java.nio.charset.StandardCharsets;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.PrivateKey;
import java.security.Provider;
import java.security.PublicKey;
import java.security.Security;
import java.security.Signature;
import java.util.Base64;

public class PKCS11SignatureService implements SignatureService {

    private final String pkcs11ConfigPath;
    private final String pin;
    private final String alias;

    public PKCS11SignatureService(String pkcs11ConfigPath, String pin, String alias) {
        if (pkcs11ConfigPath == null || pkcs11ConfigPath.isEmpty()) {
            throw new IllegalArgumentException("PKCS#11 config path is required.");
        }
        if (pin == null || pin.isEmpty()) {
            throw new IllegalArgumentException("PIN is required for PKCS#11 access.");
        }
        if (alias == null || alias.isEmpty()) {
            throw new IllegalArgumentException("Key alias is required.");
        }
        this.pkcs11ConfigPath = pkcs11ConfigPath;
        this.pin = pin;
        this.alias = alias;
    }

    @Override
    public String sign(String data, String algorithm) throws Exception {
        Provider provider = getProvider();
        try {
            Security.addProvider(provider);
            KeyStore keyStore = KeyStore.getInstance("PKCS11", provider);
            keyStore.load(null, pin.toCharArray());

            PrivateKey privateKey = (PrivateKey) keyStore.getKey(alias, null);
            if (privateKey == null) {
                throw new KeyStoreException("Key not found for alias: " + alias);
            }

            Signature signature = Signature.getInstance(algorithm, provider);
            signature.initSign(privateKey);
            signature.update(data.getBytes(StandardCharsets.UTF_8));

            byte[] signatureBytes = signature.sign();
            return Base64.getEncoder().encodeToString(signatureBytes);
        } finally {
            Security.removeProvider(provider.getName());
        }
    }

    @Override
    public boolean validate(String data, String signature, String algorithm) throws Exception {
        Provider provider = getProvider();
        try {
            Security.addProvider(provider);
            KeyStore keyStore = KeyStore.getInstance("PKCS11", provider);
            keyStore.load(null, pin.toCharArray());

            java.security.cert.Certificate cert = keyStore.getCertificate(alias);
            if (cert == null) {
                throw new KeyStoreException("Certificate not found for alias: " + alias);
            }
            PublicKey publicKey = cert.getPublicKey();

            Signature sig = Signature.getInstance(algorithm, provider);
            sig.initVerify(publicKey);
            sig.update(data.getBytes(StandardCharsets.UTF_8));

            byte[] signatureBytes = Base64.getDecoder().decode(signature);
            return sig.verify(signatureBytes);
        } finally {
            Security.removeProvider(provider.getName());
        }
    }

    private Provider getProvider() {
        return new sun.security.pkcs11.SunPKCS11(this.pkcs11ConfigPath);
    }
}