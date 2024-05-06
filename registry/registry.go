package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"serviceRegistry/services"
	"time"
)

func main() {
	registry := services.NewRegistry()
	port := os.Getenv("PORT")
	playOnRegistry(port, registry)
	idleServer := registry.AvailableServer
	errorCounter := 0
	for {
		for i, s := range idleServer {
			for wt := 1 * time.Second; wt <= 10*time.Second; wt += 1 * time.Second {
				conn, err := net.DialTimeout("tcp", s, wt)
				if conn != nil {
					_ = conn.Close()
					break
				}
				if err != nil {
					errorCounter++
				}
			}
			if errorCounter >= 10 {
				fmt.Printf("The server %s fail \n", s)
				registry.AvailableServer = append(registry.AvailableServer[:i], registry.AvailableServer[i+1:]...)
			}
			errorCounter = 0
		}
		//Update the available server list for load balancer
		if len(registry.LoadBalancerAddress) != 0 {
			idleServer = updateLoadBalancer(registry.LoadBalancerAddress, registry.AvailableServer)
		}
		time.Sleep(1 * time.Second)
	}
}

func updateLoadBalancer(loadBalancer string, updatedList []string) []string {
	var idleServer []string
	lb, err := rpc.Dial("tcp", loadBalancer)
	if err != nil {
		fmt.Printf("An error occured %s", err)
	}
	err = lb.Call("LoadBalancer.UpdateAvailableServers", updatedList, &idleServer)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err)
	}
	return idleServer
}

func playOnRegistry(port string, registry *services.Registry) {
	server := rpc.NewServer()
	err := server.RegisterName("Registry", registry)
	if err != nil {
		fmt.Printf("format of Register service is not correct: %s", err)
	}
	port = fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", port)
	go server.Accept(lis)
}
