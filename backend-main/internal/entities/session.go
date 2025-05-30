package entities

import "time"

type SessionState = int

const (
	SessionActive   = SessionState(0)
	SessionPause    = SessionState(1)
	SessionFinished = SessionState(2)
)

type Session struct {
	Id             int          `db:"id" json:"id"`
	State          SessionState `db:"state" json:"state"`
	MachineId      string       `db:"machine_id" json:"machineId"`
	WorkerId       int          `db:"worker_id" json:"workerId"`
	DatetimeStart  time.Time    `json:"datetimeStart"`
	DatetimeFinish time.Time    `json:"datetimeFinish"`
}
