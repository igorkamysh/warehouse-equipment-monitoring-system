package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/service"
	"github.com/pkg/errors"
)

func getMachineCurrentMacAddr(machine *entities.Machine, timeout time.Duration) (macAddr string, err error) {
	address := fmt.Sprintf("http://%s/%s/get_mac_addr", machine.IPAddr, machine.Id)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return "", errors.Wrap(err, "create new request")
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := struct {
		MacAddr string `json:"router_bssid"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed send to arduino current status")
	}

	return data.MacAddr, nil
}

func sendMachineCurrentState(machine *entities.Machine, timeout time.Duration) error {
	payload := []byte(fmt.Sprintf(`{"current_state": %d}`, machine.State))
	reader := bytes.NewReader(payload)

	address := fmt.Sprintf("http://%s/%s", machine.IPAddr, machine.Id)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address, reader)
	if err != nil {
		return errors.Wrap(err, "create new request")
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed send to arduino current status")
	}

	return nil
}

func canUnlockMachine(svc *service.Service, user *entities.User, _ *entities.Machine) error {
	switch user.JobPosition {
	case entities.Worker:
		unfinishedSessions, err := svc.GetUnfinishedSessionsByUserId(user.Id)
		if err != nil {
			return errors.Wrap(err, "get active sessions by userId")
		}

		if len(unfinishedSessions) != 0 {
			msg := fmt.Sprintf("user have active sessions, cnt=%d", len(unfinishedSessions))
			return errors.Wrap(err, msg)
		}
		return nil

	case entities.Admin:
		return nil

	default:
		return errors.New("user has uknown job position")
	}
}

func canUnstopMachine(svc *service.Service, user *entities.User, machine *entities.Machine) (*entities.Session, error) {
	if user.JobPosition == entities.Worker {
		pausedSessions, err := svc.GetPausedSessionsByMachineAndUser(machine.Id, user.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get paused sessions by machine and user")
		}

		if len(pausedSessions) == 0 {
			return nil, errors.New("there is no paused session")
		}

		if len(pausedSessions) > 1 {
			return nil, errors.New("there are several paused sessions with machine and user")
		}
		return &pausedSessions[0], nil
	}

	if user.JobPosition == entities.Admin {
		pausedSessions, err := svc.GetPausedSessionsByMachineID(machine.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get paused sessions by machine.Id")
		}

		if len(pausedSessions) == 0 {
			return nil, errors.New("there is no paused sessions with machine.Id")
		}

		if len(pausedSessions) > 1 {
			return nil, errors.New("there are several paused sessions with machine.Id")
		}
		return &pausedSessions[0], nil
	}

	return nil, errors.New("user has unknown job position")
}

func canLockMachine(svc *service.Service, user *entities.User, machine *entities.Machine) (*entities.Session, error) {
	if user.JobPosition == entities.Worker {
		sessions, err := svc.GetActiveSessionsByMachineAndUser(machine.Id, user.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by machine.Id and user.Id")
		}

		if len(sessions) == 0 {
			return nil, errors.New("user has no active sessions with that machine")
		}

		if len(sessions) > 1 {
			return nil, errors.New("user has several active sessions with machine")
		}

		return &sessions[0], nil
	}

	if user.JobPosition == entities.Admin {
		activeSessions, err := svc.GetActiveSessionsByMachineID(machine.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get active sessions by machine.Id")
		}
		if len(activeSessions) == 0 {
			return nil, errors.New("there is no active sessions with machine")
		}
		if len(activeSessions) > 1 {
			return nil, errors.New("there several active sessions with machine")
		}

		return &activeSessions[0], nil
	}

	return nil, errors.New("user has uknown job position")

}

func canStopMachine(svc *service.Service, user *entities.User, machine *entities.Machine) (*entities.Session, error) {
	if user.JobPosition == entities.Worker {
		sessions, err := svc.GetActiveSessionsByMachineAndUser(machine.Id, user.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by machine and user")
		}

		if len(sessions) == 0 {
			return nil, errors.Wrap(err, "there is no active sessions")
		}

		if len(sessions) > 1 {
			return nil, errors.Wrap(err, "there are several sessions with machine and user")
		}
		return &sessions[0], nil
	}

	if user.JobPosition == entities.Admin {
		sessions, err := svc.GetActiveSessionsByMachineID(machine.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by machine")
		}

		if len(sessions) == 0 {
			return nil, errors.Wrap(err, "there is no sessions with machine")
		}

		if len(sessions) > 1 {
			return nil, errors.Wrap(err, "machine has several active sessions")
		}
		return &sessions[0], nil
	}

	return nil, errors.New("unknown user's job")
}

func newQrKey() string {
	gen := rand.New(rand.NewSource(time.Now().Unix()))
	n := gen.Int()

	return strconv.Itoa(int(n))
}
