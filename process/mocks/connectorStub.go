package mocks

import "github.com/ElrondNetwork/node-monitoring/data"

// ConnectorStub implements process.Connector interface
type ConnectorStub struct {
	GetEventCalled func() (data.NotificationMessage, error)
}

// GetEvent -
func (c *ConnectorStub) GetEvent() (data.NotificationMessage, error) {
	if c.GetEventCalled != nil {
		return c.GetEventCalled()
	}

	return data.NotificationMessage{}, nil
}

// IsInterfaceNil -
func (c *ConnectorStub) IsInterfaceNil() bool {
	return c == nil
}
