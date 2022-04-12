package monitoring

// eventsHandler defines the behaviour of an events processor
type eventsHandler interface {
	Run()
	Close() error
}
