package com.github.kyriosdata.assinador;

import java.util.Arrays;

public class Main {
    public static void main(String[] args) {
        if (args.length == 0) {
            System.err.println("Comando não fornecido. Use 'sign' ou 'validate'.");
            System.exit(1);
        }

        String command = args[0];
        String[] params = Arrays.copyOfRange(args, 1, args.length);

        SignatureService service = new FakeSignatureService();

        try {
            switch (command) {
                case "sign":
                    // US-02.2: Validação de parâmetros de criação
                    if (params.length < 2) {
                        throw new IllegalArgumentException("Parâmetros insuficientes para 'sign'. Uso: sign <dados> <algoritmo>");
                    }
                    String dataToSign = params[0];
                    String signAlgorithm = params[1];
                    String signature = service.sign(dataToSign, signAlgorithm);
                    System.out.println("Assinatura (simulada): " + signature);
                    break;

                case "validate":
                    // US-02.3: Validação de parâmetros de validação
                    if (params.length < 3) {
                        throw new IllegalArgumentException("Parâmetros insuficientes para 'validate'. Uso: validate <dados> <assinatura> <algoritmo>");
                    }
                    String dataToValidate = params[0];
                    String signatureToValidate = params[1];
                    String validationAlgorithm = params[2];
                    boolean isValid = service.validate(dataToValidate, signatureToValidate, validationAlgorithm);
                    System.out.println("Resultado da validação (simulada): " + (isValid ? "Válida" : "Inválida"));
                    break;

                default:
                    System.err.println("Comando desconhecido: " + command);
                    System.exit(1);
            }
        } catch (Exception e) {
            System.err.println("Erro: " + e.getMessage());
            System.exit(1);
        }
    }
}