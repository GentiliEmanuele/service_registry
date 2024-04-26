package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"regexp"
	"serviceRegistry/services"
	"time"
)

func main() {
	registry := services.NewRegistry()
	port := os.Getenv("PORT")
	playOnRegistry(port, registry)
	for {
		for key, _ := range registry.AvailableServer {
			conn, err := net.DialTimeout("tcp", key, 8*time.Second)
			if conn != nil {
				_ = conn.Close()
			}
			if err != nil {
				fmt.Printf("The server %s fail \n", key)
				delete(registry.AvailableServer, key)
			}
		}
		//Update the available server list for load balancer
		if len(registry.LoadBalancerAddress) != 0 {
			updateLoadBalancer(registry.LoadBalancerAddress, registry.AvailableServer)
		}
		time.Sleep(1 * time.Second)
	}
}

func updateLoadBalancer(loadBalancer string, updatedList map[string]string) {
	var flag bool
	lb, err := rpc.Dial("tcp", loadBalancer)
	if err != nil {
		fmt.Printf("An error occured %s", err)
	}
	err = lb.Call("LoadBalancer.UpdateAvailableServers", updatedList, &flag)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err)
	}
}

func playOnRegistry(port string, registry *services.Registry) {
	server := rpc.NewServer()
	err := server.RegisterName("Registry", registry)
	if err != nil {
		fmt.Printf("format of Register service is not correct: %s", err)
	}
	port = fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", port)
	fmt.Printf("The Register service listening on the port %s%s\n", GetOutboundIP().String(), port)
	go server.Accept(lis)
}

func getPort() string {
	var match string
	file, err := os.Open("Dockerfile")
	if err != nil {
		fmt.Printf("Error opening DockerFile %s\n", err)
		return ""
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	//REX that match the strings in the docker file that correspond to EXPOSE
	portPattern := regexp.MustCompile(`EXPOSE\s+(\d+)`)
	//Read file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := portPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			match = matches[1]
		}
	}
	return match
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
