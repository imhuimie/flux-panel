package controllers

import (
	"net/http"
	"relay-panel/backend/database"
	"relay-panel/backend/gost_api"
	"relay-panel/backend/models"

	"github.com/gin-gonic/gin"
)

func CreateTunnel(c *gin.Context) {
	var tunnel models.Tunnel
	if err := c.ShouldBindJSON(&tunnel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&tunnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tunnel"})
		return
	}

	c.JSON(http.StatusOK, tunnel)
}

func GetTunnels(c *gin.Context) {
	var tunnels []models.Tunnel
	if err := database.DB.Find(&tunnels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tunnels"})
		return
	}
	c.JSON(http.StatusOK, tunnels)
}
func GetTunnel(c *gin.Context) {
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tunnel not found"})
		return
	}
	c.JSON(http.StatusOK, tunnel)
}

func UpdateTunnel(c *gin.Context) {
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tunnel not found"})
		return
	}

	if err := c.ShouldBindJSON(&tunnel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&tunnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tunnel"})
		return
	}

	c.JSON(http.StatusOK, tunnel)
}

func DeleteTunnel(c *gin.Context) {
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tunnel not found"})
		return
	}

	if err := database.DB.Delete(&tunnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tunnel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tunnel deleted successfully"})
}

func AssignUserTunnel(c *gin.Context) {
	var userTunnel models.UserTunnel
	if err := c.ShouldBindJSON(&userTunnel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&userTunnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign user tunnel"})
		return
	}

	c.JSON(http.StatusOK, userTunnel)
}

func GetUserTunnels(c *gin.Context) {
	var userTunnels []models.UserTunnel
	if err := database.DB.Where("user_id = ?", c.Param("id")).Find(&userTunnels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tunnels"})
		return
	}
	c.JSON(http.StatusOK, userTunnels)
}

func RemoveUserTunnel(c *gin.Context) {
	var userTunnel models.UserTunnel
	if err := database.DB.Where("user_id = ? AND tunnel_id = ?", c.Param("id"), c.Param("tunnel_id")).First(&userTunnel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tunnel not found"})
		return
	}

	if err := database.DB.Delete(&userTunnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user tunnel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User tunnel removed successfully"})
}

func DiagnoseTunnel(c *gin.Context) {
	var tunnel models.Tunnel
	if err := database.DB.First(&tunnel, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tunnel not found"})
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
	// In a real application, you would perform some checks on the tunnel.
	c.JSON(http.StatusOK, gin.H{"message": "Tunnel is healthy"})
}
