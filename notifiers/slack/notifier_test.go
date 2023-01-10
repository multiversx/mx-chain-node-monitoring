package slack_test

import (
	"testing"

	"github.com/multiversx/mx-chain-node-monitoring/config"
	"github.com/multiversx/mx-chain-node-monitoring/data"
	"github.com/multiversx/mx-chain-node-monitoring/mocks"
	"github.com/multiversx/mx-chain-node-monitoring/notifiers/slack"
	"github.com/stretchr/testify/require"
)

func createMockSlackNotifierArgs() slack.ArgsSlackNotifier {
	return slack.ArgsSlackNotifier{
		Config: &config.Slack{
			Enabled: true,
			URL:     "http://localhost",
		},
		HTTPClient: &mocks.HTTPClientStub{},
	}
}

func TestNewSlackNotifier(t *testing.T) {
	t.Parallel()

	t.Run("empty url string", func(t *testing.T) {
		t.Parallel()

		args := createMockSlackNotifierArgs()
		args.Config.URL = ""

		sn, err := slack.NewSlackNotifier(args)
		require.Nil(t, sn)
		require.Equal(t, slack.ErrInvalidSlackURL, err)
	})

	t.Run("nil http client", func(t *testing.T) {
		t.Parallel()

		args := createMockSlackNotifierArgs()
		args.HTTPClient = nil

		sn, err := slack.NewSlackNotifier(args)
		require.Nil(t, sn)
		require.Equal(t, slack.ErrNilHTTPClient, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		sn, err := slack.NewSlackNotifier(createMockSlackNotifierArgs())
		require.NotNil(t, sn)
		require.Nil(t, err)
	})
}

func TestPushMessage(t *testing.T) {
	t.Parallel()

	args := createMockSlackNotifierArgs()

	wasCalled := false
	args.HTTPClient = &mocks.HTTPClientStub{
		CallPostRestEndPointCalled: func(address, path string, data interface{}) error {
			wasCalled = true
			return nil
		},
	}

	sn, err := slack.NewSlackNotifier(args)
	require.Nil(t, err)

	err = sn.PushMessage(data.NotificationMessage{})
	require.Nil(t, err)
	require.True(t, wasCalled)
}
