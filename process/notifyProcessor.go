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
func (np *notifyProcessor) AddNotifier(notifier Notifier) {
	np.mutWorkers.Lock()
	np.workers[notifier.GetID()] = notifier
	np.mutWorkers.Unlock()
}

// PushMessage will push notification message to all registered workers
func (np *notifyProcessor) PushMessage(msg data.NotificationMessage) {
	np.mutWorkers.RLock()
	for _, worker := range np.workers {
		go worker.PushMessage(msg)
	}
	np.mutWorkers.RUnlock()
}

// IsInterfaceNil returns true if there is no value under the interface
func (np *notifyProcessor) IsInterfaceNil() bool {
	return np == nil
}
