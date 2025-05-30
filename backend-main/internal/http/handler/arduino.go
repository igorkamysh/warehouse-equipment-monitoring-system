package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) RegisterMachine(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.RegisterMachine")
	var data struct {
		MachineId string `json:"machine_id"`
		IPAddr    string `json:"ip_addr"`
	}

	if err := utils.ParseRequestData(r.Body, &data); err != nil {
		slog.Error("parse req data", op, slog.String("error", err.Error()))
		if err := utils.RespondWith400(w, "failed to parse request data"); err != nil {
			slog.Error("failed to respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	idAttr, ipAttr := slog.String("machineId", data.MachineId), slog.String("ipAddr", data.IPAddr)

	machine, err := h.service.GetMachineByID(data.MachineId)
	if err != nil {
		machine, err = h.service.InsertMachine(data.MachineId, data.IPAddr)
		if err != nil {
			slog.Error("failed to create new in machine", idAttr, ipAttr, slog.String("error", err.Error()))
			if err = utils.RespondWith400(w, "failed to create new machine"); err != nil {
				slog.Error("failed to respond with 400", op, slog.String("error", err.Error()))
			}
			return
		}

	} else {
		machine, err = h.service.UpdateMachineIPAddr(machine.Id, data.IPAddr)
		if err != nil {
			slog.Error("update machine IP", idAttr, ipAttr, slog.String("error", err.Error()))

			if err = utils.RespondWith400(w, "failed to update machine IP Address"); err != nil {
				slog.Error("failed to respond with 400", op, slog.String("error", err.Error()))
			}
			return
		}
	}

	payload := struct {
		CurrentState int `json:"current_state"`
	}{CurrentState: machine.State}

	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond Success(200) with paylod on RegisterMachine",
				slog.Any("payload", payload), idAttr, ipAttr,
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
		}
	}
}
