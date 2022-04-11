package common

type EventLevel int

const (
	InfoEvent EventLevel = iota
	AlarmEvent
	CriticalEvent
)
