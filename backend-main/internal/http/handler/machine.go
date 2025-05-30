package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetMachinesByParkingName(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.GetMachinesByParkingId")

	name := r.URL.Query().Get("name")
	parking, err := h.service.GetParkingByName(name)
	if err != nil {
		slog.Error("parking not found", slog.String("parking_name", name))
		if err := utils.RespondWith400(w, "parking not found. Missing or invalid query field name"); err != nil {
			slog.Error(
				"failed respond 400",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	machines, err := h.service.GetMachinesByParkingId(parking.Id)
	if err != nil {
		slog.Error("get machines by parking_id", op, slog.String("error", err.Error()))

		if err := utils.RespondWith400(w, "failed to get machines by parking_id"); err != nil {
			slog.Error("failed respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, machines); err != nil {
		slog.Error("failed respond with JSON", op, slog.String("error", err.Error()))
	}
}

func (h *Handler) GetAllMachines(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.GetAllMachines")

	machines, err := h.service.GetAllMachines()
	if err != nil {
		slog.Error("get all machines", op, slog.String("error", err.Error()))

		if err := utils.RespondWith400(w, "failed to get all machines"); err != nil {
			slog.Error("failed respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, machines); err != nil {
		slog.Error("failed respond with JSON", op, slog.String("error", err.Error()))
	}
}

func (h *Handler) GetMachineByID(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.GetmachineByID")

	machineId := r.URL.Query().Get("machine_id")
	machine, err := h.service.GetMachineByID(machineId)

	if err != nil {
		slog.Error("get machine from db", op, slog.String("machine_id", machineId),
			slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "failed get machine by id"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))

		}
		return
	}

	if err = utils.RespondWithJSON(w, 200, machine); err != nil {
		slog.Error("failed to respond with json with machine", op, slog.String("machine_id", machineId),
			slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "failed respond with JSON"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))

		}
	}
}
