package types

// Args Variable for Register service
type Args struct {
	IPAddress, PortNumber string
}
type Flag bool

// GetServicesInput Variable for getService service
type GetServicesInput string
type ListOfServer []string
