package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/ecol-master/sharing-wh-machines/internal/dbs/postgres"
	"github.com/ecol-master/sharing-wh-machines/internal/http/handler"
	"github.com/pkg/errors"
)

type App struct {
	server *http.Server
	cfg    *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		server: &http.Server{},
		cfg:    cfg,
	}
}

// Function will panic if can not connect to db
func (a *App) Run() error {
	db, err := postgres.New(a.cfg.Postgres)
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to postgres db"))
	}

	slog.Info("successfully connect to database")
	handler := handler.New(db, a.cfg).MakeHTTPHandler()
	slog.Info("successfully initialize http handlers")

	addr := fmt.Sprintf("%s:%d", a.cfg.App.Addr, a.cfg.App.Port)
	slog.Info("staring app", slog.String("address", addr))
	return http.ListenAndServe(addr, handler)
}
