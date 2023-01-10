package noderating

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-node-monitoring/clients"
	"github.com/multiversx/mx-chain-node-monitoring/common"
	"github.com/multiversx/mx-chain-node-monitoring/config"
	"github.com/multiversx/mx-chain-node-monitoring/data"
)

var log = logger.GetOrCreate("clients/nodeRating")

const (
	maxRating        = 100
	defaultLastValue = -1.0
)

const (
	// TODO: handle a more generic path; we should be able to provide also node's api
	nodesBLSKeyPath = "/nodes/%s"
)

// ArgsNodeRating defines the arguments needed to create a new client
type ArgsNodeRating struct {
	Client clients.HTTPClient
	Config *config.NodeRating
}

type nodeRating struct {
	httpClient clients.HTTPClient
	lastValues map[string]float64
	firstRun   bool
	config     *config.NodeRating
}

// NewNodeRatingClient creates an instance of httpClient which is a wrapper for http.Client
func NewNodeRatingClient(args ArgsNodeRating) (*nodeRating, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	lastValues := make(map[string]float64)
	for _, pubKey := range args.Config.PubKeys {
		lastValues[pubKey] = defaultLastValue
	}

	return &nodeRating{
		httpClient: args.Client,
		lastValues: lastValues,
		firstRun:   true,
		config:     args.Config,
	}, nil
}

func checkArgs(args ArgsNodeRating) error {
	if check.IfNil(args.Client) {
		return ErrNilHTTPClient
	}
	if len(args.Config.PubKeys) == 0 {
		return ErrEmptyPubKeys
	}
	if len(args.Config.ApiUrl) == 0 {
		return ErrEmptyApiUrl
	}
	if args.Config.Threshold <= 0 {
		return fmt.Errorf("%w: invalid node rating threshold, provided %.2f", common.ErrInvalidValue, args.Config.Threshold)
	}

	return nil
}

func (hcw *nodeRating) GetEvent() (data.NotificationMessage, error) {
	if hcw.firstRun == true {
		log.Info("First run. Will not trigger any event.")
		hcw.firstRun = false

		return hcw.handleFirstRun()
	}

	return hcw.handleEvents()
}

func (hcw *nodeRating) handleEvents() (data.NotificationMessage, error) {
	event := data.NotificationMessage{Level: common.InfoEvent}

	nodes, err := hcw.fetchAPINodesByBLSKey()
	if err != nil {
		return event, err
	}

	msg := ""
	for _, node := range nodes {
		if hcw.lastValues[node.Bls] != defaultLastValue {
			// TODO: evaluate simply if the rating drops under a specified threshold
			diff := node.TempRating - hcw.lastValues[node.Bls]
			if diff >= 0 {
				continue
			}

			changePercentage := math.Abs(diff) / float64(maxRating)
			changePercentage = changePercentage * 100
			if changePercentage > hcw.config.Threshold {
				event.Level = common.CriticalEvent
				nodeMsg := fmt.Sprintf(
					"NodeName: %s - TempRating decreased with %.1f percent, current value: %.2f, last value: %.2f\n",
					node.Name,
					hcw.config.Threshold,
					node.TempRating,
					hcw.lastValues[node.Bls],
				)
				msg = msg + nodeMsg
			}
		}

		hcw.lastValues[node.Bls] = node.TempRating
	}

	event.Message = msg

	return event, nil
}

func (hcw *nodeRating) handleFirstRun() (data.NotificationMessage, error) {
	nodes, err := hcw.fetchAPINodesByBLSKey()
	if err != nil {
		return data.NotificationMessage{}, err
	}

	for _, node := range nodes {
		hcw.lastValues[node.Bls] = node.TempRating
	}

	return data.NotificationMessage{Level: common.NoEvent}, nil
}

func (hcw *nodeRating) fetchAPINodesByBLSKey() ([]clients.APINode, error) {
	nodes := make([]clients.APINode, 0)

	for _, pubKey := range hcw.config.PubKeys {
		path := fmt.Sprintf(nodesBLSKeyPath, pubKey)
		responseBodyBytes, err := hcw.httpClient.CallGetRestEndPoint(hcw.config.ApiUrl, path)
		if err != nil {
			return nil, err
		}

		var response clients.APINode
		err = json.Unmarshal(responseBodyBytes, &response)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, response)
	}

	return nodes, nil
}

// GetID will return using id for client
func (hcw *nodeRating) GetID() string {
	return "NodeRating"
}

// IsInterfaceNil returns true if there is no value under the interface
func (hcw *nodeRating) IsInterfaceNil() bool {
	return hcw == nil
}
