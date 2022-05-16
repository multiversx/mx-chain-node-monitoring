package common

// EventLevel defines event level type
type EventLevel int

const (
	// NoEvent specifies when no event is triggered
	NoEvent EventLevel = iota
	// InfoEvent defines a general info event type
	InfoEvent
	// CriticalEvent defines a critical event type
	CriticalEvent
)
