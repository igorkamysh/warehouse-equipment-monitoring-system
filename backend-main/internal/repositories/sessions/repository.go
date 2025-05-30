package sessions

import (
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) InsertSession(userId int, machineId string) (*entities.Session, error) {
	var id int

	timeNow := time.Now().Unix()

	q := `INSERT INTO sessions (machine_id, worker_id, datetime_start, datetime_finish) VALUES ($1, $2, $3, $4) RETURNING id;`

	if err := r.db.QueryRowx(q, machineId, userId, timeNow, timeNow).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "insert new session and scan id")
	}

	return r.GetSessionByID(id)
}

func (r *repository) GetSessionByID(sessionId int) (*entities.Session, error) {
	var (
		session        entities.Session
		timeStartUnix  int64
		timeFinishUnix int64
	)

	q := `SELECT id, state, machine_id, worker_id, datetime_start, datetime_finish FROM sessions WHERE id = $1`
	if err := r.db.QueryRowx(q, sessionId).Scan(&session.Id, &session.State, &session.MachineId, &session.WorkerId, &timeStartUnix, &timeFinishUnix); err != nil {
		return nil, errors.Wrap(err, "get session and scan values")
	}

	session.DatetimeStart = time.Unix(timeStartUnix, 0)
	session.DatetimeFinish = time.Unix(timeFinishUnix, 0)

	return &session, nil
}

func (r *repository) GetAllSessions() ([]entities.Session, error) {
	q := `SELECT id, state, machine_id, worker_id, datetime_start, datetime_finish FROM sessions`

	return r.selectSessions(q)
}

func (r *repository) GetActiveSessionsByMachineID(machineId string) ([]entities.Session, error) {
	q := `SELECT id, state, machine_id, worker_id, datetime_start, datetime_finish FROM sessions WHERE machine_id = $1 AND state = 0`

	return r.selectSessions(q, machineId)
}

func (r *repository) GetPausedSessionsByMachineID(machineId string) ([]entities.Session, error) {
	q := `SELECT * FROM sessions WHERE machine_id = $1 AND state = $2`

	return r.selectSessions(q, machineId, entities.SessionPause)
}

func (r *repository) GetActiveSessionsByUserID(userId int) ([]entities.Session, error) {
	q := `SELECT * FROM sessions WHERE worker_id = $1 AND state = 0`

	return r.selectSessions(q, userId)
}

func (r *repository) GetPauseSessionsByUserID(userId int) ([]entities.Session, error) {
	q := `SELECT * FROM sessions WHERE worker_id = $1 AND state = $2`

	return r.selectSessions(q, userId, entities.SessionPause)
}

func (r *repository) GetUnfinishedSessionsByUserId(userId int) ([]entities.Session, error) {
	q := `SELECT * FROM sessions WHERE worker_id = $1 AND state != $2`

	return r.selectSessions(q, userId, entities.SessionFinished)
}

func (r *repository) GetActiveSessionsByMachineAndUser(machineId string, userId int) ([]entities.Session, error) {
	q := `SELECT * FROM sessions WHERE machine_id = $1 AND worker_id = $2 AND state = $3;`

	return r.selectSessions(q, machineId, userId, entities.SessionActive)
}

func (r *repository) GetPausedSessionsByMachineAndUser(machineId string, userId int) ([]entities.Session, error) {
	q := `SELECT * FROM sessions WHERE machine_id = $1 AND worker_id = $2 AND state = $3;`

	return r.selectSessions(q, machineId, userId, entities.SessionPause)
}

func (r *repository) UpdateSessionState(sessionId int, state entities.SessionState) (*entities.Session, error) {
	var id int

	q := `UPDATE sessions SET state = $1 WHERE id = $2 RETURNING id;`
	if err := r.db.QueryRowx(q, state, sessionId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "update session")
	}

	return r.GetSessionByID(id)
}

func (r *repository) PauseSession(sessionId int) (*entities.Session, error) {
	var id int

	q := `UPDATE sessions SET state = $1 WHERE id = $2 RETURNING id;`
	if err := r.db.QueryRowx(q, entities.SessionPause, sessionId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "update session")
	}

	return r.GetSessionByID(id)
}

func (r *repository) FinishSession(sessionId int) (*entities.Session, error) {
	var id int
	timeNow := time.Now().Unix()

	q := `UPDATE sessions SET state = $1, datetime_finish = $2 WHERE id = $3 RETURNING id;`
	if err := r.db.QueryRowx(q, entities.SessionFinished, timeNow, sessionId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "update session")
	}

	return r.GetSessionByID(id)
}

func (r *repository) selectSessions(q string, args ...any) ([]entities.Session, error) {
	sessions := make([]entities.Session, 0)

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select all sessions")
	}

	for rows.Next() {
		var (
			session        entities.Session
			timeStartUnix  int64
			timeFinishUnix int64
		)

		if err := rows.Scan(&session.Id, &session.State, &session.MachineId, &session.WorkerId, &timeStartUnix, &timeFinishUnix); err != nil {
			return nil, errors.Wrap(err, "get all sessions and scan values")
		}

		session.DatetimeStart = time.Unix(timeStartUnix, 0)
		session.DatetimeFinish = time.Unix(timeFinishUnix, 0)

		sessions = append(sessions, session)
	}

	return sessions, nil
}
