package process_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ElrondNetwork/node-monitoring/common"
	"github.com/ElrondNetwork/node-monitoring/data"
	"github.com/ElrondNetwork/node-monitoring/process"
	"github.com/ElrondNetwork/node-monitoring/process/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createNewEventMockArgs() process.ArgsEventsProcessor {
	return process.ArgsEventsProcessor{
		Pusher:             &mocks.PusherStub{},
		TriggerInternalSec: 1,
	}
}

func TestNewEventsProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil pusher", func(t *testing.T) {
		t.Parallel()

		args := createNewEventMockArgs()
		args.Pusher = nil

		ep, err := process.NewEventsProcessor(args)
		require.Nil(t, ep)
		assert.Equal(t, process.ErrNilPusher, err)
	})

	t.Run("wrong trigger interval", func(t *testing.T) {
		t.Parallel()

		args := createNewEventMockArgs()
		args.TriggerInternalSec = 0

		ep, err := process.NewEventsProcessor(args)
		require.Nil(t, ep)
		assert.True(t, errors.Is(err, common.ErrInvalidValue))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ep, err := process.NewEventsProcessor(createNewEventMockArgs())
		require.Nil(t, err)
		assert.NotNil(t, ep)
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	args := createNewEventMockArgs()

	numCalls := 0
	client := &mocks.ConnectorStub{
		GetEventCalled: func() (data.NotificationMessage, error) {
			numCalls++
			return data.NotificationMessage{}, nil
		},
	}

	ep, err := process.NewEventsProcessor(args)
	require.Nil(t, err)

	ep.AddClients(client)

	ep.Run()

	time.Sleep(time.Second*2 + time.Millisecond*500)

	ep.Close()

	assert.Equal(t, 2, numCalls)
}
