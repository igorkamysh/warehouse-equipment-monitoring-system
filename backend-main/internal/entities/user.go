package entities

type UserJob = string

const Worker = UserJob("worker")
const Admin = UserJob("admin")

type User struct {
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	PhoneNumber string  `db:"phone_number" json:"phoneNumber"`
	JobPosition UserJob `db:"job_position" json:"jobPosition"`
	Password    string  `db:"password" json:"password"`
}
