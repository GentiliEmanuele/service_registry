package types

// Args Variable for Register service
type Args struct {
	IPAddress, PortNumber, ServiceName string
}
type Flag bool

// GetServicesInput Variable for getService service
type GetServicesInput string
type ListOfServices map[string]string
