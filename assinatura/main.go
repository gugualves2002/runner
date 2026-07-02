package main

import (
	"log/slog"
	"os"

	"github.com/kyriosdata/runner/assinatura/cmd"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	if err := cmd.Execute(); err != nil {
		slog.Error("falha ao executar CLI", "error", err)
		os.Exit(1)
	}
}