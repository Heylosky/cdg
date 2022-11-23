package main

import (
	"fmt"
	"github.com/ComfortDelgro/controllers/sms"
	"github.com/ComfortDelgro/messageQ"
	"github.com/ComfortDelgro/models"
	"github.com/gin-gonic/gin"
)

var messageChan chan models.SMS
var outboundChan chan models.MbRc

func mqMiddleware(c *gin.Context) {
	c.Set("inBoundChan", messageChan)
	c.Set("outBoundChan", outboundChan)
}

func main() {
	// 初始化多个rabbitMQ客户端
	messageChan = make(chan models.SMS)
	for i := 0; i < 2; i++ {
		go messageQ.InitMessageQ(messageChan, i, "inboundSMS")
	}

	outboundChan = make(chan models.MbRc)
	for i := 0; i < 2; i++ {
		go messageQ.InitMbQ(outboundChan, i, "outboundSMS")
	}

	r := gin.Default()

	smsRouters := r.Group("/sms")
	smsRouters.Use(mqMiddleware)
	{
		smsRouters.POST("/v1/send", sms.SendSms)
		smsRouters.POST("/v1/receive", sms.ReceiveSms)
	}

	r.POST("/test", func(context *gin.Context) {
		message, _ := context.GetPostForm("payload")
		fmt.Println(message)
		context.JSON(200, gin.H{
			"data": message,
		})
	})

	r.Run()
}
