package main

import (
	"log/slog"
	"os"
	"path"

	"github.com/ecol-master/sharing-wh-machines/internal/app"
	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/ecol-master/sharing-wh-machines/internal/libs/csv"
	"github.com/ecol-master/sharing-wh-machines/internal/logger"
	"github.com/pkg/errors"
)

func main() {
	cfg := config.MustLoad()
	setupLoggers(cfg)

	a := app.New(cfg)
	if err := a.Run(); err != nil {
		slog.Error(err.Error())
	}
}

func setupLoggers(cfg *config.Config) {
	if err := os.Mkdir(cfg.Log.OutDir, 0750); err != nil && !os.IsExist(err) {
		panic(errors.Wrap(err, "create logs folder"))
	}

	logger.Setup(path.Join(cfg.Log.OutDir, cfg.Log.Dev))
	csv.Setup(path.Join(cfg.Log.OutDir, cfg.Log.CSV))
}
