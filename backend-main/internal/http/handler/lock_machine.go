package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/libs/csv"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) LockMachine(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.LockMachine")
	// TODO: распарсить данные для работы
	var data struct {
		MachineId string `json:"machine_id"`
	}

	err := utils.ParseRequestData(r.Body, &data)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed parse request body"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	machine, err := h.service.GetMachineByID(data.MachineId)
	if err != nil {
		slog.Error("get machine by id", op, slog.String("machine_id", data.MachineId),
			slog.String("error", err.Error()))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get machine by id",
				slog.String("data_machine_id", data.MachineId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// Получаем mac адрес от машинки
	currentMac, err := getMachineCurrentMacAddr(machine, h.cfg.MC.RequestTimeout)
	if err != nil {
		slog.Error("failed getMachineCurrentMacAddr", op, slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "can't get machine mac addr at the moment"); err != nil {
			slog.Error("failed to respond with 400 on getMachineCurrentMacAddr",
				slog.Any("machine", machine),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// Проверяем, что парковка с таким мак-адресом существует
	parking, err := h.service.GetParkingByMacAddr(currentMac)
	if err != nil {
		slog.Error("failed GetParkingByMacAddr", op, slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "can't get parking by macaddr. Parking not exists"); err != nil {
			slog.Error("failed to respond with 400 on GetParkingByMacAddr",
				slog.Any("machine", machine),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// Проверяем, что парковка может принять ещё одну машинку
	if int(parking.Capacity) <= parking.Machines && parking.Capacity != 0 {
		if err = utils.RespondWith400(w, "error while adding machine to parking. Parking machines is more or equals than capacity"); err != nil {
			slog.Error("failed to respond 400 on failed adding machine to parking",
				slog.Int("parking_id", parking.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// Проверяем, что парковка активна
	if parking.State == entities.ParkingInactive {
		if err = utils.RespondWith400(w, "error while adding machine to parking. Parking is inactive for now"); err != nil {
			slog.Error("failed to respond 400 on failed adding machine to parking",
				slog.Int("parking_id", parking.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if machine.State != entities.MachineInUse {
		if err = utils.RespondWith400(w, "machine is not in use at the moment"); err != nil {
			if err = utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 during machine is not in use",
					slog.String("machine_id", data.MachineId),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
		}
		return
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		slog.Error("get `user_id` from r.Context", op, slog.Any("context", r.Context()))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get 'use_id' value from r.Context",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	user, err := h.service.GetUserByID(int(userId))
	if err != nil {
		slog.Error("get user by id", op, slog.Int64("user_id", userId), slog.String("error", err.Error()))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get user by id",
				slog.Int64("user_id", userId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	session, err := canLockMachine(h.service, user, machine)
	if err != nil {
		slog.Error("tryLockMachine", op, slog.Int("user_id", user.Id),
			slog.String("machine_id", machine.Id), slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "user can not lock machine"); err != nil {
			slog.Error("failed to respond 500 on failed to try lock machine",
				slog.Int64("user_id", userId),
				slog.String("machine_id", machine.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// TODO: обновить данные в базе данных у машины
	machine.State = entities.MachineFree
	_, err = h.service.UpdateMachineState(machine.Id, machine.State)
	if err != nil {
		slog.Error("failed to update machine state LockMachine",
			slog.Any("machine", machine),
			slog.Int("new_state", machine.State),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed update machine state LockMachine",
				slog.String("machine_id", machine.Id),
				slog.Int("new_state", machine.State),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// TODO: Отправить данные на машину, для обработки
	if err := sendMachineCurrentState(machine, h.cfg.MC.RequestTimeout); err != nil {
		slog.Error("send machine new state", slog.String("machine_id", machine.Id),
			slog.Int("new_state", machine.State), slog.String("error", err.Error()))
		return
	}

	// TODO: завершить сессию
	session, err = h.service.FinishSession(session.Id)
	if err != nil {
		// TODO:
		slog.Error("failed to update session state UnlockMachine",
			slog.Any("machine", machine),
			slog.Int("new_state", machine.State),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed session machine state UnlockMachine",
				slog.Int("session_id", session.Id),
				slog.Int("new_state", entities.SessionFinished),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}

		return
	}

	// TODO: ответить, что добавление сессии прошло успешно
	slog.Info("session was successfully stopped", slog.Int("session_id", session.Id))

	// log to csv file information about ending of the session
	csv.Write(csv.CsvData{
		UserId:          userId,
		UserName:        user.Name,
		SessionStart:    session.DatetimeStart,
		SessionDuration: session.DatetimeFinish.Sub(session.DatetimeStart),
	})

	// Если всё хорошо - добавляем машинку на парковку
	_, err = h.service.UpdateParkingMachines(parking.Machines+1, parking.Id)
	if err != nil {
		slog.Error("parking with such id doesn't exists", slog.Int("parking_id", parking.Id))
		if err = utils.RespondWith400(w, "parking with such id doesn't exists"); err != nil {
			if err = utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 during parking with such id doesn't exists",
					slog.Int("parking_id", machine.ParkingId),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
		}
		return
	}

	// И обновляем id парковки у машинки
	_, err = h.service.UpdateMachineParkingId(machine.Id, parking.Id)
	if err != nil {
		slog.Error("failed to update machine parkingId lockMachine",
			slog.Any("machine", machine),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed update machine parkingId lockMachine",
				slog.String("machine_id", machine.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	payload := struct {
		Msg string `json:"msg"`
	}{Msg: "successfullly lock machine"}

	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		slog.Error("failed to respond with 200 on lock machine",
			slog.String("machine_id", machine.Id),
			slog.Int("user_id", user.Id),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}

}
