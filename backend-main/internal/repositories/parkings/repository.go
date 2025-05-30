package parkings

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

// Add new parking
func (r *repository) InsertParking(name, mac string, capacity entities.Capacity, state entities.ParkingState) (*entities.Parking, error) {
	var parking entities.Parking

	q := `
		INSERT INTO parkings (name, mac_addr, capacity, state)
		VALUES ($1, $2, $3, $4)
		RETURNING *;
	`
	if err := r.db.QueryRowx(q, name, mac, capacity, state).StructScan(&parking); err != nil {
		return nil, errors.Wrap(err, "inserting new parking")
	}

	return &parking, nil
}

// Get parking via Id
func (r *repository) GetParkingById(parkingId int) (*entities.Parking, error) {
	var parking entities.Parking

	q := `SELECT * FROM parkings WHERE id = $1`
	if err := r.db.Get(&parking, q, parkingId); err != nil {
		return nil, errors.Wrap(err, "get parking by id")
	}
	return &parking, nil
}

// Get parking via Name
func (r *repository) GetParkingByName(name string) (*entities.Parking, error) {
	var parking entities.Parking

	q := `SELECT * FROM parkings WHERE name = $1`
	if err := r.db.Get(&parking, q, name); err != nil {
		return nil, errors.Wrap(err, "get parking by id")
	}
	return &parking, nil
}

// Get parking via MacAddr
func (r *repository) GetParkingByMacAddr(macAddr string) (*entities.Parking, error) {
	var parking entities.Parking

	q := `SELECT * FROM parkings WHERE mac_addr = $1`
	if err := r.db.Get(&parking, q, macAddr); err != nil {
		return nil, errors.Wrap(err, "get parking by mac_addr")
	}
	return &parking, nil
}

// Get all parkings
func (r *repository) GetAllParkings() ([]entities.Parking, error) {
	parkings := make([]entities.Parking, 0)

	q := `SELECT * FROM parkings`
	if err := r.db.Select(&parkings, q); err != nil {
		return nil, errors.Wrap(err, "select all parkings")
	}
	return parkings, nil
}

// Update parking state (active, inactive)
func (r *repository) UpdateParkingState(state entities.ParkingState, parkingId int) (*entities.Parking, error) {
	var id int

	q := `UPDATE parkings SET state = $1 WHERE id = $2 RETURNING id;`
	if err := r.db.QueryRowx(q, state, parkingId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "update parking state")
	}

	return r.GetParkingById(id)
}

// Add or remove machines from parking
func (r *repository) UpdateParkingMachines(machines int, parkingId int) (*entities.Parking, error) {
	var id int

	q := `UPDATE parkings SET machines = $1 WHERE id = $2 RETURNING id;`
	if err := r.db.QueryRowx(q, machines, parkingId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "update parking machines")
	}

	return r.GetParkingById(id)
}

// Update parking capacity
func (r *repository) UpdateParkingCapacity(capacity entities.Capacity, parkingId int) (*entities.Parking, error) {
	var id int

	q := `UPDATE parkings SET capacity = $1 WHERE id = $2 RETURNING id;`
	if err := r.db.QueryRowx(q, capacity, parkingId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "update parking capacity")
	}

	return r.GetParkingById(id)
}
