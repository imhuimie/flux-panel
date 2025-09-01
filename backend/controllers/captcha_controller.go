package controllers

import (
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// TODO: In a distributed system, you would use a shared store like Redis.
var store = captcha.NewMemoryStore(10240, 5*time.Minute)

func init() {
	captcha.SetCustomStore(store)
}

func GenerateCaptcha(c *gin.Context) {
	length := captcha.DefaultLen
	captchaId := captcha.NewLen(length)
	c.JSON(http.StatusOK, gin.H{
		"captchaId": captchaId,
		"imageUrl":  "/api/captcha/" + captchaId + ".png",
	})
}

func ServeCaptcha(c *gin.Context) {
	captchaId := c.Param("captchaId")
	if captchaId == "" {
		c.String(http.StatusNotFound, "not found")
		return
	}
	// Strip the .png extension
	captchaId = captchaId[:len(captchaId)-4]

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	err := captcha.WriteImage(c.Writer, captchaId, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		if err == captcha.ErrNotFound {
			c.String(http.StatusNotFound, "not found")
		} else {
			c.String(http.StatusInternalServerError, "internal error")
		}
	}
}

func VerifyCaptcha(c *gin.Context) {
	var json struct {
		CaptchaId       string `json:"captchaId"`
		CaptchaSolution string `json:"captchaSolution"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if captcha.VerifyString(json.CaptchaId, json.CaptchaSolution) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid captcha"})
	}
}
