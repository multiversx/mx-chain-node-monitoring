package notifiers

// HTTPClient defines the behaviour of a http client
type HTTPClient interface {
	CallGetRestEndPoint(address string, path string) ([]byte, error)
	CallPostRestEndPoint(address string, path string, data interface{}) error
	IsInterfaceNil() bool
}
