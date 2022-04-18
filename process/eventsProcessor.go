package process

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/common"
)

var log = logger.GetOrCreate("process")

const minTriggerIntervalSec = 1

// ArgsEventsProcessor defines the arguments needed for events processor creation
type ArgsEventsProcessor struct {
	Pusher             Pusher
	TriggerInternalSec int
}

type eventsProcessor struct {
	clients            map[string]Connector
	mutClients         sync.RWMutex
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
		clients:            make(map[string]Connector),
		pusher:             args.Pusher,
		triggerInternalSec: args.TriggerInternalSec,
	}, nil
}

func checkArgs(args ArgsEventsProcessor) error {
	if check.IfNil(args.Pusher) {
		return ErrNilPusher
	}
	if args.TriggerInternalSec < minTriggerIntervalSec {
		return fmt.Errorf("%w: minimum trigger interval in seconds %d, provided %d", common.ErrInvalidValue, args.TriggerInternalSec, minTriggerIntervalSec)
	}

	return nil
}

func (ep *eventsProcessor) AddClients(client Connector) {
	ep.mutClients.RLock()
	ep.clients[client.GetID()] = client
	ep.mutClients.RUnlock()
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
	ep.mutClients.RLock()
	defer ep.mutClients.RUnlock()

	for id, client := range ep.clients {
		event, err := client.GetEvent()
		if err != nil {
			log.Error("failed to get event for client", "client", id, "error", err.Error())
			return
		}

		switch event.Level {
		case common.CriticalEvent:
			log.Info("Critical Event received. Will try to send event.", "clientID", id)
			ep.pusher.PushMessage(event)
		case common.InfoEvent:
			log.Info("Info event received. Will not send notification.", "clientID", id)
		case common.NoEvent:
			log.Debug("No event received. Will not send notification.", "clientID", id)
		default:
			log.Error("Invalid event level", "clientID", id)
		}
	}
}

func (ep *eventsProcessor) Close() error {
	if ep.cancelFunc != nil {
		ep.cancelFunc()
	}

	return nil
}
