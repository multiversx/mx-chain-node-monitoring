package noderating_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ElrondNetwork/node-monitoring/clients"
	"github.com/ElrondNetwork/node-monitoring/clients/mocks"
	noderating "github.com/ElrondNetwork/node-monitoring/clients/nodeRating"
	"github.com/ElrondNetwork/node-monitoring/common"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createDefaultMockArgs() noderating.ArgsNodeRating {
	return noderating.ArgsNodeRating{
		Client: &mocks.HTTPClientStub{},
		Config: &config.NodeRating{
			Threshold: 1.0,
			ApiUrl:    "http://localhost:8080",
			PubKeys:   []string{"pubk1"},
		},
	}
}

func TestNewNodeRatingClient(t *testing.T) {
	t.Parallel()

	t.Run("nil http client", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMockArgs()
		args.Client = nil

		nr, err := noderating.NewNodeRatingClient(args)
		require.Nil(t, nr)
		assert.Equal(t, noderating.ErrNilHTTPClient, err)
	})

	t.Run("no public keys provided in config", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMockArgs()
		args.Config.PubKeys = []string{}

		nr, err := noderating.NewNodeRatingClient(args)
		require.Nil(t, nr)
		assert.Equal(t, noderating.ErrEmptyPubKeys, err)
	})

	t.Run("invalid threshold value in config", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMockArgs()
		args.Config.Threshold = 0

		nr, err := noderating.NewNodeRatingClient(args)
		require.Nil(t, nr)
		assert.True(t, errors.Is(err, common.ErrInvalidValue))
	})

	t.Run("empty api url in config", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMockArgs()
		args.Config.ApiUrl = ""

		nr, err := noderating.NewNodeRatingClient(args)
		require.Nil(t, nr)
		assert.Equal(t, noderating.ErrEmptyApiUrl, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		nr, err := noderating.NewNodeRatingClient(createDefaultMockArgs())
		require.Nil(t, err)
		assert.NotNil(t, nr)
	})
}

func TestGetEvent(t *testing.T) {
	t.Parallel()

	t.Run("first run", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMockArgs()

		testAPINode := &clients.APINode{Bls: "blskey"}
		testAPINodeBytes, _ := json.Marshal(testAPINode)

		args.Client = &mocks.HTTPClientStub{
			CallGetRestEndPointCalled: func(address, path string) ([]byte, error) {
				return testAPINodeBytes, nil
			},
		}

		nr, err := noderating.NewNodeRatingClient(args)
		require.Nil(t, err)

		assert.True(t, nr.GetFirstRun())

		event, err := nr.GetEvent()
		require.Nil(t, err)

		expectedEvent := data.NotificationMessage{Level: common.NoEvent}

		assert.Equal(t, expectedEvent, event)
		assert.False(t, nr.GetFirstRun())
	})

	t.Run("after first run", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMockArgs()
		args.Config.PubKeys = []string{"blskey"}

		testAPINode := &clients.APINode{
			Bls:        "blskey",
			TempRating: 100,
		}
		testAPINodeBytes, _ := json.Marshal(testAPINode)

		testAPINode2 := &clients.APINode{
			Bls:        "blskey",
			TempRating: 90,
		}
		testAPINodeBytes2, _ := json.Marshal(testAPINode2)

		numCalls := 0
		args.Client = &mocks.HTTPClientStub{
			CallGetRestEndPointCalled: func(address, path string) ([]byte, error) {
				if numCalls == 0 {
					numCalls++
					return testAPINodeBytes, nil
				}

				return testAPINodeBytes2, nil
			},
		}

		nr, err := noderating.NewNodeRatingClient(args)
		require.Nil(t, err)

		// First Run
		event, err := nr.GetEvent()
		require.Nil(t, err)
		assert.False(t, nr.GetFirstRun())

		// Second Run
		event, err = nr.GetEvent()
		require.Nil(t, err)

		assert.Equal(t, common.CriticalEvent, event.Level)
	})
}
