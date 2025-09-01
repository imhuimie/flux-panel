package gost_api

import (
	"strconv"
	"strings"
)

func (c *GostClient) AddService(nodeID int64, name string, inPort int, limiter int, remoteAddr string, fowType int, tunnel map[string]interface{}, strategy string, interfaceName string) ([]byte, error) {
	var services []interface{}
	protocols := []string{"tcp", "udp"}
	for _, protocol := range protocols {
		service := createServiceConfig(name, inPort, limiter, remoteAddr, protocol, fowType, tunnel, strategy, interfaceName)
		services = append(services, service)
	}
	return c.Send(strconv.FormatInt(nodeID, 10), "AddService", services)
}

func (c *GostClient) UpdateService(nodeID int64, name string, inPort int, limiter int, remoteAddr string, fowType int, tunnel map[string]interface{}, strategy string, interfaceName string) ([]byte, error) {
	var services []interface{}
	protocols := []string{"tcp", "udp"}
	for _, protocol := range protocols {
		service := createServiceConfig(name, inPort, limiter, remoteAddr, protocol, fowType, tunnel, strategy, interfaceName)
		services = append(services, service)
	}
	return c.Send(strconv.FormatInt(nodeID, 10), "UpdateService", services)
}

func (c *GostClient) DeleteService(nodeID int64, name string) ([]byte, error) {
	data := make(map[string][]string)
	services := []string{name + "_tcp", name + "_udp"}
	data["services"] = services
	return c.Send(strconv.FormatInt(nodeID, 10), "DeleteService", data)
}

func createServiceConfig(name string, inPort int, limiter int, remoteAddr, protocol string, fowType int, tunnel map[string]interface{}, strategy, interfaceName string) map[string]interface{} {
	service := make(map[string]interface{})
	service["name"] = name + "_" + protocol

	if protocol == "tcp" {
		service["addr"] = tunnel["tcp_listen_addr"].(string) + ":" + strconv.Itoa(inPort)
	} else {
		service["addr"] = tunnel["udp_listen_addr"].(string) + ":" + strconv.Itoa(inPort)
	}

	if interfaceName != "" {
		metadata := make(map[string]interface{})
		metadata["interface"] = interfaceName
		service["metadata"] = metadata
	}

	if limiter != 0 {
		service["limiter"] = strconv.Itoa(limiter)
	}

	handler := createHandler(protocol, name, fowType)
	service["handler"] = handler

	listener := createListener(protocol)
	service["listener"] = listener

	if isPortForwarding(fowType) {
		forwarder := createForwarder(remoteAddr, strategy)
		service["forwarder"] = forwarder
	}

	return service
}

func createHandler(protocol, name string, fowType int) map[string]interface{} {
	handler := make(map[string]interface{})
	handler["type"] = protocol
	if isTunnelForwarding(fowType) {
		handler["chain"] = name + "_chains"
	}
	return handler
}

func createListener(protocol string) map[string]interface{} {
	listener := make(map[string]interface{})
	listener["type"] = protocol
	if protocol == "udp" {
		metadata := make(map[string]interface{})
		metadata["keepAlive"] = true
		listener["metadata"] = metadata
	}
	return listener
}

func createForwarder(remoteAddr, strategy string) map[string]interface{} {
	forwarder := make(map[string]interface{})
	var nodes []interface{}
	split := strings.Split(remoteAddr, ",")
	for i, addr := range split {
		node := make(map[string]interface{})
		node["name"] = "node_" + strconv.Itoa(i+1)
		node["addr"] = addr
		nodes = append(nodes, node)
	}

	if strategy == "" {
		strategy = "fifo"
	}
	forwarder["nodes"] = nodes

	selector := make(map[string]interface{})
	selector["strategy"] = strategy
	selector["maxFails"] = 1
	selector["failTimeout"] = "600s"
	forwarder["selector"] = selector

	return forwarder
}

func isPortForwarding(fowType int) bool {
	return fowType == 1
}

func isTunnelForwarding(fowType int) bool {
	return fowType != 1
}
