package handler

import (
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/http/middlewares"
	"github.com/ecol-master/sharing-wh-machines/internal/service"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	service *service.Service
	cfg     *config.Config
	qrKey   string
}

func New(db *sqlx.DB, cfg *config.Config) *Handler {
	return &Handler{
		service: service.New(db),
		cfg:     cfg,
		qrKey:   newQrKey(),
	}
}

func (h *Handler) MakeHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// api methods for web-application
	mux.Handle("GET /get_all_users", h.makeAdminHandler(h.GetAllUsers))
	mux.Handle("GET /get_user", h.makeAdminHandler(h.GetUserByID))

	mux.Handle("GET /get_parking_machines", h.makeAdminHandler(h.GetMachinesByParkingName))
	mux.Handle("GET /get_all_machines", h.makeAdminHandler(h.GetAllMachines))
	mux.Handle("GET /get_machine", h.makeAdminHandler(h.GetMachineByID))

	mux.Handle("GET /get_all_sessions", h.makeAdminHandler(h.GetAllSessions))
	mux.Handle("GET /get_session", h.makeAdminHandler(h.GetSessionByID))

	// New handlers for parkings
	mux.Handle("GET /get_all_parkings", h.makeAdminHandler(h.GetAllParkings))
	mux.Handle("GET /get_parking", h.makeAdminHandler(h.GetParkingById))
	mux.Handle("POST /register_parking", h.makeAdminHandler(h.RegisterParking))
	mux.Handle("PUT /update_parking_state", h.makeAdminHandler(h.UpdateParkingState))
	mux.Handle("PUT /update_parking_capacity", h.makeAdminHandler(h.UpdateParkingCapacity))
	mux.Handle("PUT /add_machine", h.makeAdminHandler(h.ManualyMoveParkingMachine))

	// auth
	mux.HandleFunc("POST /login", h.Login)

	// handlers to work with qr
	mux.Handle("GET /get_qr_key", h.makeWorkerHandler(h.GetQrKey))
	mux.Handle("POST /finish_session", h.makeWorkerHandler(h.FinishSession))

	// Lock, Unlock, Pause handler
	mux.Handle("POST /unlock_machine", h.makeWorkerHandler(h.UnlockMachine))
	mux.Handle("POST /lock_machine", h.makeWorkerHandler(h.LockMachine))
	mux.Handle("POST /stop_machine", h.makeWorkerHandler(h.StopMachine))
	mux.Handle("POST /unstop_machine", h.makeWorkerHandler(h.UnstopMachine))

	// handler to register (or make active after failed) arduino in system
	mux.Handle("POST /register_machine", http.HandlerFunc(h.RegisterMachine))

	// logging all request with LoggingMiddleware
	return middlewares.CorsEnableMiddleware(middlewares.LoggingMiddleware(mux))
}

// function making handler from RoleBasedAccess middleware with entities.Admin role
func (h *Handler) makeAdminHandler(handleFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return middlewares.RoleBasedAccess(h.cfg.Secret, entities.Admin, http.HandlerFunc(handleFunc))
}

func (h *Handler) makeWorkerHandler(handleFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return middlewares.RoleBasedAccess(h.cfg.Secret, entities.Worker, http.HandlerFunc(handleFunc))
}
