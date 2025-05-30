package entities

type ParkingState int

const ParkingInactive = ParkingState(0)
const ParkingActive = ParkingState(1)

type Capacity int

const UnlimitedCapacity = Capacity(0)

type Parking struct {
	Id       int          `db:"id" json:"id"`
	Name     string       `db:"name" json:"name"`
	MacAddr  string       `db:"mac_addr" json:"mac_addr"`
	Machines int          `db:"machines" json:"machines"`
	Capacity Capacity     `db:"capacity" json:"capacity"`
	State    ParkingState `db:"state" json:"state"`
}
