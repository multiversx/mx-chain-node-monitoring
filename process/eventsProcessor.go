package process

import (
	"context"
	"fmt"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/common"
)

var log = logger.GetOrCreate("process")

// ArgsEventsProcessor defines the arguments needed for events processor creation
type ArgsEventsProcessor struct {
	Client             Connector
	Pusher             Pusher
	TriggerInternalSec int
}

type eventsProcessor struct {
	client             Connector
	pusher             Pusher
	triggerInternalSec time.Duration
	cancelFunc         func()
}

// NewEventsProcessor will create a new events processor instance
func NewEventsProcessor(args ArgsEventsProcessor) (*eventsProcessor, error) {
	ep := &eventsProcessor{
		client: args.Client,
		pusher: args.Pusher,
	}

	return ep, nil
}

// Run will trigger the main process loop
func (ep *eventsProcessor) Run() {
	var ctx context.Context
	ctx, ep.cancelFunc = context.WithCancel(context.Background())

	ep.run(ctx)
}

func (ep *eventsProcessor) run(ctx context.Context) {
	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Debug("events processor is stopping...")
			return
		case <-timer.C:
			ep.handleEvents()
		}
	}
}

func (ep *eventsProcessor) handleEvents() {
	event, err := ep.client.GetEvent()
	if err != nil {
		fmt.Println("failed to get event")
		return
	}

	switch event.Level {
	case common.CriticalEvent:
		log.Info("Critical Event received. Will try to send event...")
		ep.pusher.PushMessage(event)
	case common.InfoEvent:
		log.Info("Info event received. Will not send notification.")
	case common.NoEvent:
		log.Debug("No event received. Will not send notification.")
	default:
		log.Error("Invalid event level")
	}
}

func (ep *eventsProcessor) Close() error {
	if ep.cancelFunc != nil {
		ep.cancelFunc()
	}

	return nil
}
