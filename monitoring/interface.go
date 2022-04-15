package monitoring

// processorHandler defines the behaviour of an events processor
type processorHandler interface {
	Run()
	Close() error
}
