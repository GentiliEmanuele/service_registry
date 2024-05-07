package services

import (
	"fmt"
	"net/rpc"
	"os"
	"serviceRegistry/types"
	"sync"
)

// Registry : Create a struct that maintain the state of the server

type Registry struct {
	AvailableServer     []string
	LoadBalancerAddress string
	mapMutex            sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		AvailableServer:     make([]string, 0),
		LoadBalancerAddress: "",
		mapMutex:            sync.RWMutex{},
	}
}

func (s *Registry) Register(args *types.Args, ret *types.Flag) error {
	s.mapMutex.Lock()
	address := fmt.Sprintf("%s:%s", args.IPAddress, args.PortNumber) //All server are identified by IP and port number pairs
	//Memorize the server
	s.AvailableServer = append(s.AvailableServer, address)
	//Send the new server to load balancer
	s.LoadBalancerAddress = os.Getenv("LOAD_BALANCER")
	loadBalancer, err := rpc.Dial("tcp", s.LoadBalancerAddress)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err)
	}
	err = loadBalancer.Call("LoadBalancer.AddNewServer", address, &s.AvailableServer)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err)
	}
	*ret = true
	s.mapMutex.Unlock()
	return nil
}
