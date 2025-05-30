package logger

import (
	"log/slog"
	"os"

	"github.com/pkg/errors"
)

func Setup(filepath string) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(errors.Wrap(err, "open dev logs file"))
	}
	defer file.Close()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)
}
