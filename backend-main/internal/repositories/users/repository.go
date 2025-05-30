package users

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

func (r *repository) InsertUser(name, phoneNumber, jobPosition, password string) (*entities.User, error) {
	var user entities.User
	q := `
		INSERT INTO users(name, phone_number, job_position, password)
		VALUES ($1, $2, $3, $4)
		RETURNING;
	`
	if err := r.db.QueryRowx(q, name, phoneNumber, jobPosition, password).StructScan(&user); err != nil {
		return nil, errors.Wrap(err, "inserting new user")
	}
	return &user, nil
}

func (r *repository) GetAllUsers() ([]entities.User, error) {
	users := make([]entities.User, 0)

	q := `SELECT * FROM users`
	err := r.db.Select(&users, q)

	if err != nil {
		return nil, errors.Wrap(err, "select all users")
	}
	return users, nil
}

func (r *repository) GetUserByID(userId int) (*entities.User, error) {
	var u entities.User

	q := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&u, q, userId)
	if err != nil {
		return nil, errors.Wrap(err, "get user by id")
	}
	return &u, err
}

func (r *repository) GetUserByPhoneNumber(phoneNumber string) (*entities.User, error) {
	var u entities.User

	q := `SELECT * FROM users WHERE phone_number = $1`
	err := r.db.Get(&u, q, phoneNumber)
	if err != nil {
		return nil, errors.Wrap(err, "get user by phone number")
	}
	return &u, err
}
