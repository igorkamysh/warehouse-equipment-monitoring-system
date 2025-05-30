package entities

type MachineState = int

const MachineFree = MachineState(0)
const MachineStop = MachineState(1)
const MachineInUse = MachineState(2)

type Machine struct {
	Id        string       `db:"id" json:"id"`
	State     MachineState `db:"state" json:"state"`
	ParkingId int          `db:"parking_id" json:"parking_id"`
	Voltage   int          `db:"voltage" json:"voltage"`
	IPAddr    string       `db:"ip_addr" json:"ipAddr"`
}
