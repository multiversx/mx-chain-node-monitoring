package mocks

import "github.com/ElrondNetwork/node-monitoring/data"

// NotifierStub implements process.Notifier interface
type NotifierStub struct {
	PushMessageCalled func(msg data.NotificationMessage) error
	GetIDCalled       func() string
}

// PushMessage -
func (n *NotifierStub) PushMessage(msg data.NotificationMessage) error {
	if n.PushMessageCalled != nil {
		return n.PushMessageCalled(msg)
	}

	return nil
}

// GetID -
func (n *NotifierStub) GetID() string {
	if n.GetIDCalled != nil {
		return n.GetIDCalled()
	}

	return "ID"
}
