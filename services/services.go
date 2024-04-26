package services

import (
	"fmt"
	"sync"
)

// Registry : Create a struct that maintain the state of the server

type Registry struct {
	AvailableServer     map[string]string
	LoadBalancerAddress string
	mapMutex            sync.RWMutex
}

// Args Variable for Register service
type Args struct {
	IPAddress, PortNumber, ServiceName string
}
type Flag bool

// Address Variable for getService service
type Address string
type ListOfServices map[string]string

func NewRegistry() *Registry {
	return &Registry{
		AvailableServer:     make(map[string]string),
		LoadBalancerAddress: "",
		mapMutex:            sync.RWMutex{},
	}
}

func (s *Registry) Register(args *Args, ret *Flag) error {
	s.mapMutex.Lock()
	key := fmt.Sprintf("%s:%s", args.IPAddress, args.PortNumber) //All server are identified by IP and port number pairs
	value := args.ServiceName
	s.AvailableServer[key] = value
	s.mapMutex.Unlock()
	return nil
}

func (s *Registry) GetServices(loadBalancerIP Address, ret *ListOfServices) error {
	//Service registry save the load balancer address for update it when a sever crush
	s.LoadBalancerAddress = string(loadBalancerIP)
	fmt.Printf("Saved loadBalancerIP: %s\n", s.LoadBalancerAddress)
	*ret = s.AvailableServer
	return nil
}
