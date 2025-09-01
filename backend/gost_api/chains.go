package gost_api

import "strconv"

func (c *GostClient) AddChains(nodeID int64, name, remoteAddr, protocol, interfaceName string) ([]byte, error) {
	data := createChainData(name, remoteAddr, protocol, interfaceName)
	return c.Send(strconv.FormatInt(nodeID, 10), "AddChains", data)
}

func (c *GostClient) UpdateChains(nodeID int64, name, remoteAddr, protocol, interfaceName string) ([]byte, error) {
	data := createChainData(name, remoteAddr, protocol, interfaceName)
	req := make(map[string]interface{})
	req["chain"] = name + "_chains"
	req["data"] = data
	return c.Send(strconv.FormatInt(nodeID, 10), "UpdateChains", req)
}

func (c *GostClient) DeleteChains(nodeID int64, name string) ([]byte, error) {
	data := make(map[string]interface{})
	data["chain"] = name + "_chains"
	return c.Send(strconv.FormatInt(nodeID, 10), "DeleteChains", data)
}

func createChainData(name, remoteAddr, protocol, interfaceName string) map[string]interface{} {
	dialer := make(map[string]interface{})
	dialer["type"] = protocol
	if protocol == "quic" {
		metadata := make(map[string]interface{})
		metadata["keepAlive"] = true
		metadata["ttl"] = "10s"
		dialer["metadata"] = metadata
	}

	connector := make(map[string]interface{})
	connector["type"] = "relay"

	node := make(map[string]interface{})
	node["name"] = "node-" + name
	node["addr"] = remoteAddr
	node["connector"] = connector
	node["dialer"] = dialer

	if interfaceName != "" {
		node["interface"] = interfaceName
	}

	nodes := []interface{}{node}

	hop := make(map[string]interface{})
	hop["name"] = "hop-" + name
	hop["nodes"] = nodes

	hops := []interface{}{hop}

	data := make(map[string]interface{})
	data["name"] = name + "_chains"
	data["hops"] = hops

	return data
}
