package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/libs/csv"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetQrKey(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.GetQrKey")

	payload := struct {
		QrKey   string `json:"qr_key"`
		LocalIp string `json:"local_ip"`
	}{
		QrKey:   h.qrKey,
		LocalIp: h.cfg.App.MachineAddr,
	}

	if err := utils.RespondWithJSON(w, 200, payload); err != nil {
		slog.Error("failed respond with JSON", op, slog.String("error", err.Error()))
	}
}

func (h *Handler) FinishSession(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.FinishSession")

	data := struct {
		Key         string `json:"key"`
		ParkingName string `json:"parking_name"`
	}{}

	err := utils.ParseRequestData(r.Body, &data)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed parse request body"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	if data.Key != h.qrKey {
		slog.Error("data.Key does not match with handler's QrKey", op, slog.String("Key", data.Key))

		if err = utils.RespondWith400(w, "data.Key does not match with handler's QrKey"); err != nil {
			slog.Error("failed to respond 400 on failed data.Key does not match with handler's QrKey",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
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
		slog.Error("failed to get user by userId", op, slog.String("error", err.Error()), slog.Int("userId", int(userId)))
		if err = utils.RespondWith400(w, "failed to get user by userId"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	sessions, err := h.service.GetActiveSessionsByUserID(int(userId))
	if err != nil || sessions == nil {
		slog.Error("failed to get sessions by userId", op, slog.String("error", err.Error()), slog.Int("userId", int(userId)))
		if err = utils.RespondWith400(w, "failed to get sessions by userId"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	for _, sess := range sessions {
		machine, err := h.service.GetMachineByID(sess.MachineId)
		if err != nil {
			slog.Error("failed to get machine by session.MachineId", op, slog.String("error", err.Error()), slog.Int("userId", int(userId)), slog.Any("session", sess))
			if err = utils.RespondWith400(w, "failed to get sessions by session.MachineId"); err != nil {
				slog.Error("failed respond with 400", slog.String("error", err.Error()))
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
		parkingByMac, err := h.service.GetParkingByMacAddr(currentMac)
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

		// Проверяем, что парковка с таким именем существует
		parkingByName, err := h.service.GetParkingByName(data.ParkingName)
		if err != nil {
			slog.Error("can't get parking by name", op, slog.String("parking_name", data.ParkingName), slog.String("error", err.Error()))

			if err = utils.RespondWith400(w, "user can not lock machine"); err != nil {
				slog.Error("failed to respond 500 on can't get parking by name",
					slog.String("parking_name", data.ParkingName),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
			return
		}

		// Проверяем, что id парковок совпадают
		if parkingByMac.Id != parkingByName.Id {
			slog.Error("user trying to end session using qr from other parking place. Move machine to the qr-code's parking", op, slog.Int("user_id", user.Id))

			if err = utils.RespondWith400(w, "user trying to end session using qr from other parking place"); err != nil {
				slog.Error("failed to respond with 400 on user trying to end session using qr from other parking place",
					slog.Any("parkingByName", parkingByName),
					slog.Any("parkingByMac", parkingByMac),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
			return
		}

		// Проверяем, что парковка может принять ещё одну машинку
		if int(parkingByMac.Capacity) <= parkingByMac.Machines && parkingByMac.Capacity != 0 {
			if err = utils.RespondWith400(w, "error while adding machine to parking. Parking machines is more or equals than capacity"); err != nil {
				slog.Error("failed to respond 400 on failed adding machine to parking",
					slog.Int("parking_id", parkingByMac.Id),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
			return
		}

		// Проверяем, что парковка активна
		if parkingByMac.State == entities.ParkingInactive {
			if err = utils.RespondWith400(w, "error while adding machine to parking. Parking is inactive for now"); err != nil {
				slog.Error("failed to respond 400 on failed adding machine to parking",
					slog.Int("parking_id", parkingByMac.Id),
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

		if machine.State != entities.MachineInUse {
			if err = utils.RespondWith400(w, "machine is not in use at the moment"); err != nil {
				if err = utils.RespondWith500(w); err != nil {
					slog.Error("failed to respond with 500 during machine is not in use",
						slog.String("machine_id", machine.Id),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.String("error", err.Error()),
					)
				}
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
		_, err = h.service.UpdateParkingMachines(parkingByMac.Machines+1, parkingByMac.Id)
		if err != nil {
			slog.Error("parking with such id doesn't exists", slog.Int("parking_id", parkingByMac.Id))
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
		_, err = h.service.UpdateMachineParkingId(machine.Id, parkingByMac.Id)
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
	}

	h.qrKey = newQrKey()

	payload := struct {
		Msg string `json:"msg"`
	}{Msg: "successfullly lock machine"}

	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		slog.Error("failed to respond with 200 on lock machine",
			slog.Int("user_id", user.Id),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}
}
