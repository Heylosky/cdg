package main

import (
	"github.com/ComfortDelgro/controllers/sms"
	"github.com/ComfortDelgro/logger"
	"github.com/ComfortDelgro/messageQ"
	"github.com/ComfortDelgro/middlewares"
	"github.com/ComfortDelgro/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var messageChan chan models.SMS
var outboundChan chan models.MbRc

func mqMiddleware(c *gin.Context) {
	c.Set("inBoundChan", messageChan)
	c.Set("outBoundChan", outboundChan)
}

func main() {
	logger.InitLogger()

	// 初始化多个rabbitMQ客户端
	messageChan = make(chan models.SMS)
	for i := 0; i < 2; i++ {
		go messageQ.InitMessageQ(messageChan, i, "inboundSMS")
	}
	outboundChan = make(chan models.MbRc)
	for i := 0; i < 2; i++ {
		go messageQ.InitMbQ(outboundChan, i, "outboundSMS")
	}

	//启动http服务端，并加载日志中间件
	r := gin.New()
	r.Use(logger.GinLogger, logger.GinRecovery(false))

	smsRouters := r.Group("/sms")
	smsRouters.Use(middlewares.TokenVerify, mqMiddleware)
	{
		smsRouters.POST("/v1/send", sms.SendSms)
		smsRouters.POST("/v1/receive", sms.ReceiveSms)
	}

	userRouters := r.Group("/user")
	{
		userRouters.GET("/tokens")
	}

	r.GET("/healthz", func(context *gin.Context) {
		context.String(200, "ok")
	})

	zap.L().Info("CDG server starting.")
	r.Run("0.0.0.0:8080")
}
