package process

import (
	"sync"

	"github.com/ElrondNetwork/node-monitoring/data"
)

type notifyProcessor struct {
	workers    map[string]Notifier
	mutWorkers sync.RWMutex
}

// NewNotifyProcessor will create a new notify processor instance
func NewNotifyProcessor() *notifyProcessor {
	return &notifyProcessor{
		workers: make(map[string]Notifier),
	}
}

// AddNotifier will add a notifier instance to workers list
func (bp *notifyProcessor) AddNotifier(notifier Notifier) {
	bp.mutWorkers.RLock()
	bp.workers[notifier.GetID()] = notifier
	bp.mutWorkers.RUnlock()
}

// PushMessage will push notification message to all registered workers
func (bp *notifyProcessor) PushMessage(msg data.NotificationMessage) {
	bp.mutWorkers.RLock()
	for _, worker := range bp.workers {
		go worker.PushMessage(msg)
	}
	bp.mutWorkers.RUnlock()
}
