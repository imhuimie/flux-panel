package controllers

import (
	"net/http"
	"relay-panel/backend/database"
	"relay-panel/backend/models"

	"github.com/gin-gonic/gin"
)

func CreateNode(c *gin.Context) {
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}

	c.JSON(http.StatusOK, node)
}

func GetNodes(c *gin.Context) {
	var nodes []models.Node
	if err := database.DB.Find(&nodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nodes"})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func GetNode(c *gin.Context) {
	var node models.Node
	if err := database.DB.First(&node, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}
	c.JSON(http.StatusOK, node)
}

func UpdateNode(c *gin.Context) {
	var node models.Node
	if err := database.DB.First(&node, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}

	c.JSON(http.StatusOK, node)
}

func DeleteNode(c *gin.Context) {
	var node models.Node
	if err := database.DB.First(&node, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	if err := database.DB.Delete(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node deleted successfully"})
}

func GetInstallCommand(c *gin.Context) {
	var node models.Node
	if err := database.DB.First(&node, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	installCommand := `
#!/bin/bash
set -e

# update and install dependencies
apt-get update
apt-get install -y wget curl

# download and install gost
wget https://github.com/ginuerzh/gost/releases/download/v2.11.5/gost-linux-amd64-2.11.5.gz
gunzip gost-linux-amd64-2.11.5.gz
mv gost-linux-amd64-2.11.5 /usr/local/bin/gost
chmod +x /usr/local/bin/gost

# create gost config
cat > /etc/gost.json <<EOF
{
	   "Debug": true,
	   "Retries": 0,
	   "API": {
	       "addr": "` + node.ApiHost + `:` + node.ApiPort + `",
	       "path": "/api",
	       "accesslog": true,
	       "auth": {
	           "username": "` + node.ApiUsername + `",
	           "password": "` + node.ApiPassword + `"
	       }
	   }
}
EOF

# create gost service
cat > /etc/systemd/system/gost.service <<EOF
[Unit]
Description=GO Simple Tunnel
After=network.target
Wants=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/gost -C /etc/gost.json
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# start gost service
systemctl daemon-reload
systemctl enable gost
systemctl start gost
`
	c.JSON(http.StatusOK, gin.H{"install_command": installCommand})
}
