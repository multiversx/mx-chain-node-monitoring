package mocks

import "github.com/ElrondNetwork/node-monitoring/data"

// ConnectorStub implements process.Connector interface
type ConnectorStub struct {
	GetEventCalled func() (data.NotificationMessage, error)
	GetIDCalled    func() string
}

// GetEvent -
func (c *ConnectorStub) GetEvent() (data.NotificationMessage, error) {
	if c.GetEventCalled != nil {
		return c.GetEventCalled()
	}

	return data.NotificationMessage{}, nil
}

// GetID -
func (c *ConnectorStub) GetID() string {
	if c.GetIDCalled != nil {
		return c.GetIDCalled()
	}

	return ""
}

// IsInterfaceNil -
func (c *ConnectorStub) IsInterfaceNil() bool {
	return c == nil
}
