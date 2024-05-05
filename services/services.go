package services

import (
	"fmt"
	"serviceRegistry/types"
	"sync"
)

// Registry : Create a struct that maintain the state of the server

type Registry struct {
	AvailableServer     map[string]string
	LoadBalancerAddress string
	mapMutex            sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		AvailableServer:     make(map[string]string),
		LoadBalancerAddress: "",
		mapMutex:            sync.RWMutex{},
	}
}

func (s *Registry) Register(args *types.Args, ret *types.Flag) error {
	s.mapMutex.Lock()
	key := fmt.Sprintf("%s:%s", args.IPAddress, args.PortNumber) //All server are identified by IP and port number pairs
	value := args.ServiceName
	s.AvailableServer[key] = value
	s.mapMutex.Unlock()
	return nil
}

func (s *Registry) GetServices(loadBalancerIP types.GetServicesInput, ret *types.ListOfServices) error {
	//Service registry save the load balancer address for update it when a sever crush
	s.LoadBalancerAddress = string(loadBalancerIP)
	*ret = s.AvailableServer
	return nil
}
