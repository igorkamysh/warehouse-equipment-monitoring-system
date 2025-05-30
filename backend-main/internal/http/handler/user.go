package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()

	if err != nil {
		slog.Error(
			"failed get all users",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed respond with error", slog.Int("status", 500))
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, users); err != nil {
		slog.Error("failed to respond with json with users",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed respond with error", slog.Int("status", 500))
		}
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		slog.Error("`user_id` query is not integer")
		if err := utils.RespondWith400(w, "user_id should be integer"); err != nil {
			slog.Error(
				"failed respond 400",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	user, err := h.service.GetUserByID(id)

	if err != nil {
		slog.Error("user not found", slog.Int("user_id", id))
		if err := utils.RespondWith400(w, "user not found"); err != nil {
			slog.Error(
				"failed respond 400",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, user); err != nil {
		slog.Error("failed to respond with json with user",
			slog.Int("user_id", id),
			slog.String("error", err.Error()),
		)
		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during get user by id",
				slog.Int("user_id", id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}
}
