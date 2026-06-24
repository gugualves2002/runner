package com.kyriosdata.assinador;

import org.springframework.web.bind.annotation.*;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

@RestController
@RequestMapping("/api/v1")
public class SignController {
    private final SignatureService signatureService = new FakeSignatureService();
    private final Gson gson = new GsonBuilder().setPrettyPrinting().create();

    @PostMapping("/sign")
    public SignatureResponse sign(
            @RequestParam(value = "payload") String payload,
            @RequestParam(value = "key-alias") String keyAlias) {
        try {
            String signature = signatureService.sign(payload, keyAlias);
            return new SignatureResponse(signature, "success", "Assinatura criada com sucesso");
        } catch (SignatureException e) {
            return new SignatureResponse(null, "error", e.getMessage());
        }
    }

    @PostMapping("/validate")
    public SignatureResponse validate(
            @RequestParam(value = "payload") String payload,
            @RequestParam(value = "signature") String signature,
            @RequestParam(value = "key-alias") String keyAlias) {
        try {
            boolean isValid = signatureService.validate(payload, signature, keyAlias);
            if (isValid) {
                return new SignatureResponse(null, "valid", "Assinatura válida");
            } else {
                return new SignatureResponse(null, "invalid", "Assinatura inválida");
            }
        } catch (SignatureException e) {
            return new SignatureResponse(null, "error", e.getMessage());
        }
    }

    @GetMapping("/health")
    public SignatureResponse health() {
        return new SignatureResponse(null, "ok", "Servidor em operação");
    }
}
