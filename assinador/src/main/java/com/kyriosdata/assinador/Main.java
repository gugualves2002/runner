package com.kyriosdata.assinador;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

import java.util.Scanner;

/**
 * Classe principal da aplicação Assinador.
 * Gerencia requisições de assinatura digital e validação.
 */
public class Main {

    private static final SignatureService signatureService = new FakeSignatureService();
    private static final Gson gson = new GsonBuilder().setPrettyPrinting().create();

    public static void main(String[] args) {
        if (args.length == 0) {
            printUsage();
            System.exit(1);
        }

        try {
            String command = args[0];

            switch (command) {
                case "sign":
                    handleSign(args);
                    break;
                case "validate":
                    handleValidate(args);
                    break;
                case "version":
                    System.out.println("assinador 0.1.0");
                    break;
                case "--help":
                case "-h":
                    printHelp();
                    break;
                default:
                    System.err.println("Comando desconhecido: " + command);
                    printUsage();
                    System.exit(1);
            }
        } catch (Exception e) {
            System.err.println("Erro: " + e.getMessage());
            System.exit(1);
        }
    }

    private static void handleSign(String[] args) throws SignatureException {
        String payload = null;
        String keyAlias = null;

        // Parse de argumentos simples: sign --payload "..." --key-alias "..."
        for (int i = 1; i < args.length; i++) {
            if ("--payload".equals(args[i]) && i + 1 < args.length) {
                payload = args[i + 1];
                i++;
            } else if ("--key-alias".equals(args[i]) && i + 1 < args.length) {
                keyAlias = args[i + 1];
                i++;
            }
        }

        if (payload == null || payload.isEmpty()) {
            System.err.println("Erro: --payload é obrigatório");
            System.exit(1);
        }
        if (keyAlias == null || keyAlias.isEmpty()) {
            System.err.println("Erro: --key-alias é obrigatório");
            System.exit(1);
        }

        String signature = signatureService.sign(payload, keyAlias);
        SignatureResponse response = new SignatureResponse(signature, "success", "Assinatura criada com sucesso");

        System.out.println(gson.toJson(response));
    }

    private static void handleValidate(String[] args) throws SignatureException {
        String payload = null;
        String signature = null;
        String keyAlias = null;

        // Parse de argumentos simples
        for (int i = 1; i < args.length; i++) {
            if ("--payload".equals(args[i]) && i + 1 < args.length) {
                payload = args[i + 1];
                i++;
            } else if ("--signature".equals(args[i]) && i + 1 < args.length) {
                signature = args[i + 1];
                i++;
            } else if ("--key-alias".equals(args[i]) && i + 1 < args.length) {
                keyAlias = args[i + 1];
                i++;
            }
        }

        if (payload == null || payload.isEmpty()) {
            System.err.println("Erro: --payload é obrigatório");
            System.exit(1);
        }
        if (signature == null || signature.isEmpty()) {
            System.err.println("Erro: --signature é obrigatório");
            System.exit(1);
        }
        if (keyAlias == null || keyAlias.isEmpty()) {
            System.err.println("Erro: --key-alias é obrigatório");
            System.exit(1);
        }

        boolean isValid = signatureService.validate(payload, signature, keyAlias);
        String status = isValid ? "valid" : "invalid";
        String message = isValid ? "Assinatura válida" : "Assinatura inválida";
        SignatureResponse response = new SignatureResponse(null, status, message);

        System.out.println(gson.toJson(response));
    }

    private static void printUsage() {
        System.err.println("Uso: assinador <comando> [opções]");
        System.err.println("Use 'assinador --help' para mais informações");
    }

    private static void printHelp() {
        System.out.println("Assinador - Simulação de Assinatura Digital");
        System.out.println();
        System.out.println("Uso: assinador <comando> [opções]");
        System.out.println();
        System.out.println("Comandos:");
        System.out.println("  sign              Criar uma assinatura digital");
        System.out.println("  validate          Validar uma assinatura digital");
        System.out.println("  version           Exibir versão");
        System.out.println("  --help, -h        Exibir esta ajuda");
        System.out.println();
        System.out.println("Opções para 'sign':");
        System.out.println("  --payload <texto>      Conteúdo a ser assinado");
        System.out.println("  --key-alias <alias>    Alias da chave privada");
        System.out.println();
        System.out.println("Opções para 'validate':");
        System.out.println("  --payload <texto>      Conteúdo original");
        System.out.println("  --signature <sig>      Assinatura em Base64");
        System.out.println("  --key-alias <alias>    Alias da chave pública");
        System.out.println();
        System.out.println("Exemplos:");
        System.out.println("  assinador sign --payload 'Olá Mundo' --key-alias minha-chave");
        System.out.println("  assinador validate --payload 'Olá Mundo' --signature '...' --key-alias minha-chave");
    }
}
