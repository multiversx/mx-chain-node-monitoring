package process_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/node-monitoring/data"
	"github.com/ElrondNetwork/node-monitoring/process"
	"github.com/ElrondNetwork/node-monitoring/process/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyProcessor(t *testing.T) {
	t.Parallel()

	np := process.NewNotifyProcessor()
	require.False(t, np.IsInterfaceNil())

	notification := data.NotificationMessage{
		Message: "message",
		Level:   1,
	}

	wg := sync.WaitGroup{}

	numCalls := uint32(0)
	notifier := &mocks.NotifierStub{
		PushMessageCalled: func(msg data.NotificationMessage) error {
			assert.Equal(t, notification, msg)
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
			return nil
		},
		GetIDCalled: func() string {
			return "ID"
		},
	}

	np.AddNotifier(notifier)

	wg.Add(1)

	np.PushMessage(notification)

	wg.Wait()

	assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}
