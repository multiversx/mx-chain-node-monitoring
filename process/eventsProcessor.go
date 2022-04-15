package process

import (
	"context"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/common"
)

var log = logger.GetOrCreate("process")

const minTriggerIntervalSec = 1

// ArgsEventsProcessor defines the arguments needed for events processor creation
type ArgsEventsProcessor struct {
	Client             Connector
	Pusher             Pusher
	TriggerInternalSec int
}

type eventsProcessor struct {
	client             Connector
	pusher             Pusher
	triggerInternalSec int
	cancelFunc         func()
}

// NewEventsProcessor will create a new events processor instance
func NewEventsProcessor(args ArgsEventsProcessor) (*eventsProcessor, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &eventsProcessor{
		client:             args.Client,
		pusher:             args.Pusher,
		triggerInternalSec: args.TriggerInternalSec,
	}, nil
}

func checkArgs(args ArgsEventsProcessor) error {
	if check.IfNil(args.Client) {
		return ErrNilClient
	}
	if check.IfNil(args.Pusher) {
		return ErrNilPusher
	}
	if args.TriggerInternalSec < minTriggerIntervalSec {
		return fmt.Errorf("%w: minimum trigger interval in seconds %d, provided %d", common.ErrInvalidValue, args.TriggerInternalSec, minTriggerIntervalSec)
	}

	return nil
}

// Run will trigger the main process loop
func (ep *eventsProcessor) Run() {
	var ctx context.Context
	ctx, ep.cancelFunc = context.WithCancel(context.Background())

	go ep.run(ctx)
}

func (ep *eventsProcessor) run(ctx context.Context) {
	timer := time.NewTicker(time.Second * time.Duration(ep.triggerInternalSec))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("events processor is stopping...")
			return
		case <-timer.C:
			ep.handleEvents()
		}
	}
}

func (ep *eventsProcessor) handleEvents() {
	event, err := ep.client.GetEvent()
	if err != nil {
		log.Error("failed to get event", "error", err.Error())
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
