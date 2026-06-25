package com.github.kyriosdata.assinador;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.kyriosdata.assinador.model.*;
import io.javalin.Javalin;
import io.javalin.http.Context;

import static io.javalin.apibuilder.ApiBuilder.*;

public class WebServer {

    private final SignatureService signatureService;
    private final ObjectMapper objectMapper = new ObjectMapper();

    public WebServer(SignatureService signatureService) {
        this.signatureService = signatureService;
    }

    public void start() {
        Javalin app = Javalin.create().start(7070);

        app.routes(() -> {
            path("/api", () -> {
                post("/sign", this::handleSign);
                post("/validate", this::handleValidate);
            });
        });

        // Tratamento de exceções para retornar JSON
        app.exception(Exception.class, (e, ctx) -> {
            ctx.status(400);
            ctx.json(new ErrorResponse(e.getMessage()));
        });
    }

    private void handleSign(Context ctx) throws Exception {
        SignatureRequest req = objectMapper.readValue(ctx.body(), SignatureRequest.class);
        SignatureService service = getServiceForRequest(req.getPkcs11ConfigPath(), req.getPin(), req.getAlias());

        String signature = service.sign(req.getData(), req.getAlgorithm());
        ctx.status(200).json(new SignatureResponse(signature));
    }

    private void handleValidate(Context ctx) throws Exception {
        ValidationRequest req = objectMapper.readValue(ctx.body(), ValidationRequest.class);
        SignatureService service = getServiceForRequest(req.getPkcs11ConfigPath(), req.getPin(), req.getAlias());

        boolean isValid = service.validate(req.getData(), req.getSignature(), req.getAlgorithm());
        ctx.status(200).json(new ValidationResponse(isValid));
    }

    private SignatureService getServiceForRequest(String pkcs11ConfigPath, String pin, String alias) {
        // If PKCS#11 config is provided, create a specific service for this request.
        // Otherwise, use the default (fake) service.
        if (pkcs11ConfigPath != null && !pkcs11ConfigPath.isEmpty()) {
            return new PKCS11SignatureService(pkcs11ConfigPath, pin, alias);
        }
        return this.signatureService;
    }
}