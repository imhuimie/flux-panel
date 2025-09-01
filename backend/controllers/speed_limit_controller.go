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

func CreateSpeedLimit(c *gin.Context) {
	var speedLimit models.SpeedLimit
	if err := c.ShouldBindJSON(&speedLimit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&speedLimit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create speed limit"})
		return
	}

	// After creating the speed limit in db, update all related nodes
	var nodes []models.Node
	database.DB.Find(&nodes) // Get all nodes

	for _, node := range nodes {
		gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
		if err := gostClient.Connect(); err != nil {
			log.Printf("Failed to connect to gost service on node %s: %v", node.Name, err)
			continue
		}
		limiterName, err := strconv.ParseInt(speedLimit.Name, 10, 64)
		if err != nil {
			log.Printf("Failed to parse speed limit name to int64: %v", err)
			continue
		}
		_, err = gostClient.AddLimiters(int64(node.ID), limiterName, speedLimit.Speed)
		if err != nil {
			log.Printf("Failed to add limiter to gost service on node %s: %v", node.Name, err)
		}
	}

	c.JSON(http.StatusOK, speedLimit)
}

func GetSpeedLimits(c *gin.Context) {
	var speedLimits []models.SpeedLimit
	if err := database.DB.Find(&speedLimits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get speed limits"})
		return
	}
	c.JSON(http.StatusOK, speedLimits)
}

func GetSpeedLimit(c *gin.Context) {
	var speedLimit models.SpeedLimit
	if err := database.DB.First(&speedLimit, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Speed limit not found"})
		return
	}
	c.JSON(http.StatusOK, speedLimit)
}

func UpdateSpeedLimit(c *gin.Context) {
	var speedLimit models.SpeedLimit
	if err := database.DB.First(&speedLimit, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Speed limit not found"})
		return
	}

	if err := c.ShouldBindJSON(&speedLimit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&speedLimit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update speed limit"})
		return
	}

	// After updating the speed limit in db, update all related nodes
	var nodes []models.Node
	database.DB.Find(&nodes) // Get all nodes

	for _, node := range nodes {
		gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
		if err := gostClient.Connect(); err != nil {
			log.Printf("Failed to connect to gost service on node %s: %v", node.Name, err)
			continue
		}
		limiterName, err := strconv.ParseInt(speedLimit.Name, 10, 64)
		if err != nil {
			log.Printf("Failed to parse speed limit name to int64: %v", err)
			continue
		}
		_, err = gostClient.UpdateLimiters(int64(node.ID), limiterName, speedLimit.Speed)
		if err != nil {
			log.Printf("Failed to update limiter on gost service on node %s: %v", node.Name, err)
		}
	}

	c.JSON(http.StatusOK, speedLimit)
}

func DeleteSpeedLimit(c *gin.Context) {
	var speedLimit models.SpeedLimit
	if err := database.DB.First(&speedLimit, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Speed limit not found"})
		return
	}

	// Before deleting the speed limit in db, update all related nodes
	var nodes []models.Node
	database.DB.Find(&nodes) // Get all nodes

	for _, node := range nodes {
		gostClient := gost_api.NewGostClient("ws://" + node.ApiHost + ":" + node.ApiPort + "/api")
		if err := gostClient.Connect(); err != nil {
			log.Printf("Failed to connect to gost service on node %s: %v", node.Name, err)
			continue
		}
		limiterName, err := strconv.ParseInt(speedLimit.Name, 10, 64)
		if err != nil {
			log.Printf("Failed to parse speed limit name to int64: %v", err)
			continue
		}
		_, err = gostClient.DeleteLimiters(int64(node.ID), limiterName)
		if err != nil {
			log.Printf("Failed to delete limiter on gost service on node %s: %v", node.Name, err)
		}
	}

	if err := database.DB.Delete(&speedLimit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete speed limit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Speed limit deleted successfully"})
}
