package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/common"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/data"
)

// TODO: refactor to separate http client wrapper from event logic

var log = logger.GetOrCreate("client/httpClient")

const (
	firstRunValue    float64 = -1
	maxRating                = 10000000
	minReqTimeoutSec         = 1
)

const (
	baseURL = "https://api.elrond.com"

	nodesIdentifierPath = "/nodes?identity=%s"
)

type httpClientWrapper struct {
	httpClient    *http.Client
	mutHTTPClient sync.RWMutex
	lastValue     float64
	config        *config.NodeRating
}

// HTTPClientWrapperArgs defines the arguments needed to create a new client
type HTTPClientWrapperArgs struct {
	ReqTimeoutSec int
	Config        *config.NodeRating
}

// NewHTTPClientWrapper creates an instance of httpClient which is a wrapper for http.Client
func NewHTTPClientWrapper(args HTTPClientWrapperArgs) (*httpClientWrapper, error) {
	if args.ReqTimeoutSec < minReqTimeoutSec {
		return nil, fmt.Errorf("%w, provided: %d, minimum: %d", common.ErrInvalidValue, args.ReqTimeoutSec, minReqTimeoutSec)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Duration(args.ReqTimeoutSec) * time.Second

	return &httpClientWrapper{
		httpClient: httpClient,
		lastValue:  firstRunValue,
		config:     args.Config,
	}, nil
}

func (hcw *httpClientWrapper) GetEvent() (data.NotificationMessage, error) {
	event := data.NotificationMessage{}

	identity := hcw.config.Identity
	percentageThreshold := hcw.config.Threshold

	if hcw.lastValue == firstRunValue {
		return event, nil
	}

	path := fmt.Sprintf(nodesIdentifierPath, identity)
	nodes, err := hcw.callGetRestEndPoint(baseURL, path)
	if err != nil {
		return event, err
	}

	msg := ""
	for _, node := range nodes {
		diff := math.Abs(node.TempRating - hcw.lastValue)
		changePercentage := diff / float64(maxRating)
		if changePercentage > percentageThreshold {
			event.Level = common.CriticalEvent
			nodeMsg := fmt.Sprintf("TempRating lower than threshold, current value: %f, threshold value: %f\n", node.TempRating, percentageThreshold)
			msg = msg + nodeMsg
		}
	}

	event.Message = msg

	return event, nil
}

// callGetRestEndPoint calls an external end point
func (hcw *httpClientWrapper) callGetRestEndPoint(
	address string,
	path string,
) ([]APINode, error) {
	fmt.Println(address)
	fmt.Println(path)
	req, err := http.NewRequest("GET", address+path, nil)
	if err != nil {
		return nil, err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Requesting data from api>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := hcw.httpClient.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return nil, err
		}

		return nil, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		if errNotCritical != nil {
			log.Warn("base process GET: close body", "error", errNotCritical.Error())
		}
	}()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response []APINode
	err = json.Unmarshal(responseBodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func isTimeoutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}
