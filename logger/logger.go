package logger

type Status int

const (
	StatusPending Status = iota
	StatusApproved
	StatusRejected
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusApproved:
		return "APPROVED"
	case StatusRejected:
		return "REJECTED"
	default:
		return "UNKNOWN"
	}
}
