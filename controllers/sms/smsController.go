package sms

import (
	"fmt"
	"github.com/ComfortDelgro/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SendSms(c *gin.Context) {
	var sms models.SMS
	if err := c.ShouldBindJSON(&sms); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":   "",
			"status": "sms format error",
		})
		return
	}

	zap.L().Info("receive messages: " + sms.Payload)

	channel, exist := c.Get("inBoundChan")
	if exist != true {
		zap.L().Error("no channel found")
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":   "",
			"status": "internal server error",
		})
	} else {
		ch1, ok := channel.(chan models.SMS)
		if ok != true {
		}
		ch1 <- sms
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   "",
		"status": "ok",
	})
}

func ReceiveSms(c *gin.Context) {
	var sms models.MbRc
	if err := c.ShouldBind(&sms); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	fmt.Println("receive messages: " + sms.Payload)

	channel, exist := c.Get("outBoundChan")
	if exist != true {
		fmt.Println("no channel found")
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":   "",
			"status": "internal server error",
		})
	} else {
		ch1, ok := channel.(chan models.MbRc)
		if ok != true {
		}
		ch1 <- sms
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   "",
		"status": "ok",
	})
}
