package common

// EventLevel defines event level type
type EventLevel int

const (
	// InfoEvent defines a general info event type
	InfoEvent EventLevel = iota
	// CriticalEvent defines a critical event type
	CriticalEvent
)
