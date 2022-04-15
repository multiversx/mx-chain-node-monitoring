package mocks

// HTTPClientStub implements HTTPClient interface
type HTTPClientStub struct {
	CallGetRestEndPointCalled func(address string, path string) ([]byte, error)
}

// CallGetRestEndPoint -
func (hcs *HTTPClientStub) CallGetRestEndPoint(address string, path string) ([]byte, error) {
	if hcs.CallGetRestEndPointCalled != nil {
		return hcs.CallGetRestEndPointCalled(address, path)
	}

	return nil, nil
}

// IsInterfaceNil -
func (hcs *HTTPClientStub) IsInterfaceNil() bool {
	return hcs == nil
}
