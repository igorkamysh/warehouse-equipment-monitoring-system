package machines

import (
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

func (r *repository) InsertMachine(machineId, ipAddr string) (*entities.Machine, error) {
	var newMachine entities.Machine

	q := `INSERT INTO machines (id, ip_addr) VALUES ($1, $2);`
	if _, err := r.db.Queryx(q, machineId, ipAddr); err != nil {
		return nil, errors.Wrap(err, "insert machine")
	}

	q = `SELECT * FROM machines WHERE id = $1;`
	if err := r.db.Get(&newMachine, q, machineId); err != nil {
		return nil, errors.Wrap(err, "select inserted machine")
	}

	return &newMachine, nil
}

func (r *repository) GetMachineByID(machineId string) (*entities.Machine, error) {
	var m entities.Machine

	q := `SELECT * FROM machines WHERE id = $1`

	if err := r.db.Get(&m, q, machineId); err != nil {
		return nil, errors.Wrap(err, "select machine by id")
	}
	return &m, nil
}

func (r *repository) GetAllMachines() ([]entities.Machine, error) {
	machines := make([]entities.Machine, 0)

	q := `SELECT * FROM machines`
	if err := r.db.Select(&machines, q); err != nil {
		return nil, errors.Wrap(err, "get all machines")
	}
	return machines, nil
}

func (r *repository) UpdateMachineIPAddr(machineId, ipAddr string) (*entities.Machine, error) {
	var machine entities.Machine

	q := `
		UPDATE machines SET ip_addr = $1 WHERE id = $2
		RETURNING *;
	`
	if err := r.db.QueryRowx(q, ipAddr, machineId).StructScan(&machine); err != nil {
		return nil, errors.Wrap(err, "failed to update machine's ipAddr")
	}
	return &machine, nil
}

func (r *repository) UpdateMachineState(machineId string, state entities.MachineState) (*entities.Machine, error) {
	var machine entities.Machine

	q := `
		UPDATE machines SET state = $1 WHERE id = $2
		RETURNING *;
	`
	if err := r.db.QueryRowx(q, state, machineId).StructScan(&machine); err != nil {
		return nil, errors.Wrap(err, "failed to update machine's ipAddr")
	}
	return &machine, nil
}

// New method to update machines parking place (parkingId)
func (r *repository) UpdateMachineParkingId(machineId string, parkingId int) (*entities.Machine, error) {
	var machine entities.Machine

	q := `
		UPDATE machines SET parking_id = $1 WHERE id = $2
		RETURNING *;
	`
	if err := r.db.QueryRowx(q, parkingId, machineId).StructScan(&machine); err != nil {
		return nil, errors.Wrap(err, "failed to update machine's parking_id")
	}
	return &machine, nil
}

// New method to select machines for each parking
func (r *repository) GetMachinesByParkingId(parkingId int) ([]entities.Machine, error) {
	machines := make([]entities.Machine, 0)

	q := `SELECT * FROM machines WHERE parking_id = $1`
	if err := r.db.Select(&machines, q, parkingId); err != nil {
		return nil, errors.Wrap(err, "get all machines")
	}
	return machines, nil
}
