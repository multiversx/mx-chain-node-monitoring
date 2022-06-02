package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/common"
)

var log = logger.GetOrCreate("clients/httpClient")

const (
	minReqTimeoutSec = 1
)

type httpClientWrapper struct {
	httpClient *http.Client
}

// HTTPClientWrapperArgs defines the arguments needed to create a new client
type HTTPClientWrapperArgs struct {
	ReqTimeoutSec int
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
	}, nil
}

// TODO: set address+path as one variable

// CallGetRestEndPoint calls an external end point
func (hcw *httpClientWrapper) CallGetRestEndPoint(
	address string,
	path string,
) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, address+path, nil)
	if err != nil {
		return nil, err
	}

	userAgent := "Elrond Node Monitoring / 1.0.0 <Requesting data from api>"
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

	return ioutil.ReadAll(resp.Body)
}

// CallPostRestEndPoint calls an external end point
func (hcw *httpClientWrapper) CallPostRestEndPoint(
	address string,
	path string,
	data interface{},
) error {
	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, address+path, bytes.NewReader(buff))
	if err != nil {
		return err
	}

	userAgent := "Elrond Node Monitoring / 1.0.0 <Requesting data from api>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := hcw.httpClient.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return err
		}

		return err
	}

	errNotCritical := resp.Body.Close()
	if errNotCritical != nil {
		log.Warn("base process GET: close body", "error", errNotCritical.Error())
	}

	return nil
}

func isTimeoutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (hcw *httpClientWrapper) IsInterfaceNil() bool {
	return hcw == nil
}
