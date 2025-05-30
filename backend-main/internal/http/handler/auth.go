package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.Login")

	var data struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	if err := utils.ParseRequestData(r.Body, &data); err != nil {
		slog.Error("parse req data", op, slog.String("error", err.Error()))
		if err := utils.RespondWith400(w, "failed to parse request data"); err != nil {
			slog.Error("failed to respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	user, err := h.service.GetUserByPhoneNumber(data.PhoneNumber)
	if err != nil {
		slog.Error("get user by phone number", op, slog.String("phone_number", data.PhoneNumber),
			slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "user with such phone number doesn't exists"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	if data.Password != user.Password {
		if err = utils.RespondWith400(w, "user password is not correct"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	token, err := h.service.GenerateToken(*user, h.cfg.Secret, h.cfg.TokenTTL)
	if err != nil {
		slog.Error("failed to generate JWT token", slog.Any("user", user), slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed to generate JWT token"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))
		}
	}

	response := struct {
		Token string `json:"token"`
	}{Token: token}

	if err := utils.RespondWithJSON(w, 200, response); err != nil {
		if err = utils.RespondWith400(w, "failed to respond successfully with JSON"); err != nil {
			slog.Error("failed to respond with JSON with JWT token", slog.String("error", err.Error()))
		}
	}
}
