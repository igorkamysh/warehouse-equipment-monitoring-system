package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.GetAllSessions()
	if err != nil {
		slog.Error("failed to get all sessions", slog.String("error", err.Error()))
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 in error get all sessions",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}
	if err = utils.RespondWithJSON(w, 200, sessions); err != nil {
		slog.Error("failed to respond with json with sessions", slog.String("error", err.Error()))
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during get all sessions",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (h *Handler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session_id")
	id, err := strconv.Atoi(sessionId)
	if err != nil {
		if err = utils.RespondWith400(w, "session_id should be integer"); err != nil {
			slog.Error("failed to respond with 400",
				slog.String("msg_respond", "session_id should be integer"),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	session, err := h.service.GetSessionByID(id)
	if err != nil {
		slog.Error("failed to get session from db",
			slog.Int("session_id", id),
			slog.String("error", err.Error()),
		)
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)

		}
		return
	}

	if err = utils.RespondWithJSON(w, 200, session); err != nil {
		slog.Error("failed to respond with json with session",
			slog.Int("session_id", id),
			slog.String("error", err.Error()),
		)
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}
}
