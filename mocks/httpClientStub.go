package mocks

// HTTPClientStub implements HTTPClient interface
type HTTPClientStub struct {
	CallGetRestEndPointCalled  func(address string, path string) ([]byte, error)
	CallPostRestEndPointCalled func(address string, path string, data interface{}) error
}

// CallGetRestEndPoint -
func (hcs *HTTPClientStub) CallGetRestEndPoint(address string, path string) ([]byte, error) {
	if hcs.CallGetRestEndPointCalled != nil {
		return hcs.CallGetRestEndPointCalled(address, path)
	}

	return nil, nil
}

// CallPostRestEndPoint -
func (hcs *HTTPClientStub) CallPostRestEndPoint(address string, path string, data interface{}) error {
	if hcs.CallPostRestEndPointCalled != nil {
		return hcs.CallPostRestEndPointCalled(address, path, data)
	}

	return nil
}

// IsInterfaceNil -
func (hcs *HTTPClientStub) IsInterfaceNil() bool {
	return hcs == nil
}
