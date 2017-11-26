package models

type Status string

const (
	Pending  Status = "pending"
	Running         = "running"
	Complete        = "complete"
	Failed          = "failed"
	Unknown         = "unknown"
)
