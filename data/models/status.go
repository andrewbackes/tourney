package models

import (
	"encoding/json"
	"strings"
)

type Status int

const (
	Pending Status = iota
	Running
	Complete
	Failed
	Unknown
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Complete:
		return "Complete"
	case Failed:
		return "Failed"
	}
	return "None"
}
func (s *Status) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	switch strings.ToLower(str) {
	default:
		*s = Unknown
	case "running":
		*s = Running
	case "complete":
		*s = Complete
	case "pending":
		*s = Pending
	case "failed":
		*s = Failed
	}
	return nil
}

func (s *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal((*s).String())
}
