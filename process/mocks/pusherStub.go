package mocks

import "github.com/multiversx/mx-chain-node-monitoring/data"

// PusherStub implements process.Pusher
type PusherStub struct {
	PushMessageCalled func(msg data.NotificationMessage)
}

// PushMessage -
func (p *PusherStub) PushMessage(msg data.NotificationMessage) {
	if p.PushMessageCalled != nil {
		p.PushMessageCalled(msg)
	}
}

// IsInterfaceNil -
func (p *PusherStub) IsInterfaceNil() bool {
	return p == nil
}
