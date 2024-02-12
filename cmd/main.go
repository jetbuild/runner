package main

import (
	"log/slog"
	"os"

	"github.com/jetbuild/runner/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	var c config.Config
	if err := c.Load(); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := c.Flow.Run(); err != nil {
		slog.Error("failed to run flow", "error", err)
		os.Exit(1)
	}
}
