package process

import (
	"sync"

	"github.com/ElrondNetwork/node-monitoring/data"
)

type ArgsNotifyProcessor struct {
}

type notifyProcessor struct {
	workers    map[string]Notifier
	mutWorkers sync.RWMutex
}

func NewNotifyProcessor(args ArgsNotifyProcessor) (*notifyProcessor, error) {
	bp := &notifyProcessor{
		workers: make(map[string]Notifier),
	}

	return bp, nil
}

func (bp *notifyProcessor) AddNotifier(notifier Notifier) {
	bp.mutWorkers.RLock()
	bp.workers[notifier.GetID()] = notifier
	bp.mutWorkers.RUnlock()
}

func (bp *notifyProcessor) PushMessage(msg data.NotificationMessage) {
	bp.mutWorkers.RLock()
	for _, worker := range bp.workers {
		go worker.PushMessage(msg)
	}
	bp.mutWorkers.RUnlock()
}
