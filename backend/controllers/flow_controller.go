package controllers

import (
	"encoding/json"
	"net/http"
	"relay-panel/backend/database"
	"relay-panel/backend/gost_api"
	"relay-panel/backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TODO: Replace with a proper client initialization and management mechanism
var GostClients = make(map[uint]*gost_api.GostClient)

type FlowData struct {
	U int64  `json:"u"`
	D int64  `json:"d"`
	N string `json:"n"`
}

func UploadFlowData(c *gin.Context) {
	secret := c.Query("secret")
	if secret == "" {
		c.String(http.StatusBadRequest, "secret is required")
		return
	}

	var node models.Node
	if err := database.DB.Where("secret = ?", secret).First(&node).Error; err != nil {
		c.String(http.StatusForbidden, "invalid secret")
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "failed to read body")
		return
	}

	// TODO: Add decryption logic here if needed

	var flowData FlowData
	if err := json.Unmarshal(body, &flowData); err != nil {
		c.String(http.StatusBadRequest, "failed to parse flow data")
		return
	}

	if flowData.N == "web_api" {
		c.String(http.StatusOK, "ok")
		return
	}

	ids := strings.Split(flowData.N, "_")
	if len(ids) != 3 {
		c.String(http.StatusBadRequest, "invalid service name")
		return
	}
	forwardIdStr := ids[0]
	userIdStr := ids[1]
	userTunnelIdStr := ids[2]

	forwardId, _ := strconv.ParseUint(forwardIdStr, 10, 64)
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)
	userTunnelId, _ := strconv.ParseUint(userTunnelIdStr, 10, 64)

	// Apply traffic ratio
	var forward models.Forward
	if err := database.DB.Preload("Tunnel").First(&forward, forwardId).Error; err == nil {
		ratio := forward.Tunnel.TrafficRatio
		if ratio > 0 {
			flowData.U = int64(float64(flowData.U) * ratio)
			flowData.D = int64(float64(flowData.D) * ratio)
		}
	} else {
		// Forward not found, probably deleted
		c.String(http.StatusOK, "ok")
		return
	}

	// Update flows
	updateForwardFlow(uint(forwardId), flowData.D, flowData.U)
	updateUserFlow(uint(userId), flowData.D, flowData.U)
	updateUserTunnelFlow(uint(userTunnelId), flowData.D, flowData.U)

	// Check limits
	checkUserLimits(uint(userId), forward.Tunnel.NodeID)
	checkUserTunnelLimits(uint(userTunnelId), forward.Tunnel.NodeID)

	c.String(http.StatusOK, "ok")
}

func updateForwardFlow(forwardId uint, inFlow int64, outFlow int64) {
	database.DB.Model(&models.Forward{}).Where("id = ?", forwardId).Updates(map[string]interface{}{
		"in_flow":  gorm.Expr("in_flow + ?", inFlow),
		"out_flow": gorm.Expr("out_flow + ?", outFlow),
	})
}

func updateUserFlow(userId uint, inFlow int64, outFlow int64) {
	database.DB.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"in_flow":  gorm.Expr("in_flow + ?", inFlow),
		"out_flow": gorm.Expr("out_flow + ?", outFlow),
	})
}

func updateUserTunnelFlow(userTunnelId uint, inFlow int64, outFlow int64) {
	if userTunnelId == 0 {
		return
	}
	database.DB.Model(&models.UserTunnel{}).Where("id = ?", userTunnelId).Updates(map[string]interface{}{
		"in_flow":  gorm.Expr("in_flow + ?", inFlow),
		"out_flow": gorm.Expr("out_flow + ?", outFlow),
	})
}

func checkUserLimits(userId uint, nodeId uint) {
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		return
	}

	// Check traffic limit
	if user.Flow > 0 && (user.InFlow+user.OutFlow) >= user.Flow {
		pauseAllUserServices(userId, nodeId)
		return
	}

	// Check expiration time
	if user.ExpTime > 0 && user.ExpTime <= time.Now().Unix() {
		pauseAllUserServices(userId, nodeId)
		return
	}

	// Check status
	if user.Status != 1 {
		pauseAllUserServices(userId, nodeId)
	}
}

func checkUserTunnelLimits(userTunnelId uint, nodeId uint) {
	if userTunnelId == 0 {
		return
	}
	var userTunnel models.UserTunnel
	if err := database.DB.First(&userTunnel, userTunnelId).Error; err != nil {
		return
	}

	// Check traffic limit
	if userTunnel.Flow > 0 && (userTunnel.InFlow+userTunnel.OutFlow) >= userTunnel.Flow {
		pauseUserServiceByTunnel(userTunnel.UserID, userTunnel.TunnelID, nodeId)
		return
	}

	// Check expiration time
	if userTunnel.ExpTime > 0 && userTunnel.ExpTime <= time.Now().Unix() {
		pauseUserServiceByTunnel(userTunnel.UserID, userTunnel.TunnelID, nodeId)
		return
	}

	// Check status
	if userTunnel.Status != 1 {
		pauseUserServiceByTunnel(userTunnel.UserID, userTunnel.TunnelID, nodeId)
	}
}

func pauseAllUserServices(userId uint, nodeId uint) {
	var forwards []models.Forward
	database.DB.Where("user_id = ?", userId).Find(&forwards)
	for _, f := range forwards {
		if client, ok := GostClients[nodeId]; ok {
			client.DeleteService(int64(nodeId), f.Name)
		}
		database.DB.Model(&f).Update("status", "paused")
	}
}

func pauseUserServiceByTunnel(userId uint, tunnelId uint, nodeId uint) {
	var forwards []models.Forward
	database.DB.Where("user_id = ? AND tunnel_id = ?", userId, tunnelId).Find(&forwards)
	for _, f := range forwards {
		if client, ok := GostClients[nodeId]; ok {
			client.DeleteService(int64(nodeId), f.Name)
		}
		database.DB.Model(&f).Update("status", "paused")
	}
}
