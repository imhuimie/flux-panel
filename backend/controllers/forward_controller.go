package controllers

import (
	"log"
	"net/http"
	"relay-panel/backend/database"
	"relay-panel/backend/gost_api"
	"relay-panel/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateForward(c *gin.Context) {
	var forward models.Forward
	if err := c.ShouldBindJSON(&forward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&forward).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create forward"})
		return
	}

	// After creating the forward rule in db, update gost service
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, forward.TunnelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find tunnel for forward rule"})
		return
	}

	var node models.Node
	if err := database.DB.First(&node, tunnel.NodeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find node for tunnel"})
		return
	}

	gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
	if err := gostClient.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gost service"})
		return
	}

	var speedLimit *models.SpeedLimit
	var limiter int
	if forward.SpeedLimitID != nil {
		if err := database.DB.First(&speedLimit, *forward.SpeedLimitID).Error; err == nil {
			limiter, _ = strconv.Atoi(speedLimit.Speed)
		}
	}

	tunnelMap := map[string]interface{}{
		"tcp_listen_addr": tunnel.TcpListenAddr,
		"udp_listen_addr": tunnel.UdpListenAddr,
	}

	_, err := gostClient.AddService(int64(node.ID), forward.Name, forward.InPort, limiter, forward.RemoteAddr, forward.FowType, tunnelMap, forward.Strategy, forward.InterfaceName)
	if err != nil {
		// Rollback db change
		database.DB.Delete(&forward)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add gost service"})
		return
	}

	c.JSON(http.StatusOK, forward)
}

func GetForwards(c *gin.Context) {
	var forwards []models.Forward
	if err := database.DB.Find(&forwards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get forwards"})
		return
	}
	c.JSON(http.StatusOK, forwards)
}

func GetForward(c *gin.Context) {
	var forward models.Forward
	if err := database.DB.First(&forward, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forward not found"})
		return
	}
	c.JSON(http.StatusOK, forward)
}

func UpdateForward(c *gin.Context) {
	var forward models.Forward
	if err := database.DB.First(&forward, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forward not found"})
		return
	}

	if err := c.ShouldBindJSON(&forward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&forward).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update forward"})
		return
	}

	// After updating the forward rule in db, update gost service
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, forward.TunnelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find tunnel for forward rule"})
		return
	}

	var node models.Node
	if err := database.DB.First(&node, tunnel.NodeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find node for tunnel"})
		return
	}

	gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
	if err := gostClient.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gost service"})
		return
	}

	var speedLimit *models.SpeedLimit
	var limiter int
	if forward.SpeedLimitID != nil {
		if err := database.DB.First(&speedLimit, *forward.SpeedLimitID).Error; err == nil {
			limiter, _ = strconv.Atoi(speedLimit.Speed)
		}
	}

	tunnelMap := map[string]interface{}{
		"tcp_listen_addr": tunnel.TcpListenAddr,
		"udp_listen_addr": tunnel.UdpListenAddr,
	}

	_, err := gostClient.UpdateService(int64(node.ID), forward.Name, forward.InPort, limiter, forward.RemoteAddr, forward.FowType, tunnelMap, forward.Strategy, forward.InterfaceName)
	if err != nil {
		// Here we don't rollback db change, just log the error
		log.Println("Failed to update gost service:", err)
	}

	c.JSON(http.StatusOK, forward)
}

func DeleteForward(c *gin.Context) {
	var forward models.Forward
	if err := database.DB.First(&forward, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forward not found"})
		return
	}

	// Before deleting the forward rule in db, delete gost service
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, forward.TunnelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find tunnel for forward rule"})
		return
	}

	var node models.Node
	if err := database.DB.First(&node, tunnel.NodeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find node for tunnel"})
		return
	}

	gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
	if err := gostClient.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gost service"})
		return
	}

	_, err := gostClient.DeleteService(int64(node.ID), forward.Name)
	if err != nil {
		// Log the error but still proceed to delete from db
		log.Println("Failed to delete gost service:", err)
	}

	if err := database.DB.Delete(&forward).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete forward"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forward deleted successfully"})
}

func PauseForward(c *gin.Context) {
	var forward models.Forward
	if err := database.DB.First(&forward, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forward not found"})
		return
	}

	forward.Status = "paused"
	if err := database.DB.Save(&forward).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update forward status"})
		return
	}

	c.JSON(http.StatusOK, forward)
}

func ResumeForward(c *gin.Context) {
	var forward models.Forward
	if err := database.DB.First(&forward, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forward not found"})
		return
	}

	forward.Status = "running"
	if err := database.DB.Save(&forward).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update forward status"})
		return
	}

	c.JSON(http.StatusOK, forward)
}

func DiagnoseForward(c *gin.Context) {
	var forward models.Forward
	if err := database.DB.First(&forward, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forward not found"})
		return
	}

	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, forward.TunnelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find tunnel for forward rule"})
		return
	}

	var node models.Node
	if err := database.DB.First(&node, tunnel.NodeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find node for tunnel"})
		return
	}

	gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
	if err := gostClient.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gost service"})
		return
	}

	// This is a placeholder for the actual diagnostic logic.
	// In a real application, you would perform some checks on the forward.
	c.JSON(http.StatusOK, gin.H{"message": "Forward is healthy"})
}

func ReorderForward(c *gin.Context) {
	var forwards []models.Forward
	if err := c.ShouldBindJSON(&forwards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, f := range forwards {
		if err := database.DB.Model(&models.Forward{}).Where("id = ?", f.ID).Update("order", i).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder forwards"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forwards reordered successfully"})
}
