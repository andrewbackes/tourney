package models

type Status int

const (
	Pending Status = iota
	Running
	Complete
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Complete:
		return "Complete"
	}
	return "None"
}
