package clients

// HTTPClient defines the behaviour of a http client
type HTTPClient interface {
	CallGetRestEndPoint(address string, path string) ([]byte, error)
	IsInterfaceNil() bool
}
