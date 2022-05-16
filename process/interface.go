package process

import "github.com/ElrondNetwork/node-monitoring/data"

// Notifier defines the behaviour of a notifier instance
type Notifier interface {
	PushMessage(msg data.NotificationMessage) error
	GetID() string
}

// Connector defines the behaviour of a client connector which will fetch the event
type Connector interface {
	GetEvent() (data.NotificationMessage, error)
	GetID() string
	IsInterfaceNil() bool
}

// Pusher defines the behaviour of a push notification instance
type Pusher interface {
	PushMessage(msg data.NotificationMessage)
	IsInterfaceNil() bool
}
