package gost_api

import "strconv"

func (c *GostClient) AddLimiters(nodeID int64, name int64, speed string) ([]byte, error) {
	data := createLimiterData(name, speed)
	return c.Send(strconv.FormatInt(nodeID, 10), "AddLimiters", data)
}

func (c *GostClient) UpdateLimiters(nodeID int64, name int64, speed string) ([]byte, error) {
	data := createLimiterData(name, speed)
	req := make(map[string]interface{})
	req["limiter"] = strconv.FormatInt(name, 10)
	req["data"] = data
	return c.Send(strconv.FormatInt(nodeID, 10), "UpdateLimiters", req)
}

func (c *GostClient) DeleteLimiters(nodeID int64, name int64) ([]byte, error) {
	req := make(map[string]interface{})
	req["limiter"] = strconv.FormatInt(name, 10)
	return c.Send(strconv.FormatInt(nodeID, 10), "DeleteLimiters", req)
}

func createLimiterData(name int64, speed string) map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = strconv.FormatInt(name, 10)
	limits := []string{"$ " + speed + "MB " + speed + "MB"}
	data["limits"] = limits
	return data
}
