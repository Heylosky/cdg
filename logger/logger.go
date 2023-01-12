package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	mode        = "dev"              //模式
	filename    = "./logs/cdg.log"   //日志存放路径
	level       = zapcore.DebugLevel //日志级别
	max_size    = 200                //最大存储大小，MB
	max_age     = 30                 //最大存储时间days
	max_backups = 10                 //备份数量
)

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	//使用lumberjack分割归档日志
	lumberLackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}
	return zapcore.AddSync(lumberLackLogger)
}

func getEncoder() zapcore.Encoder {
	//使用一份官方预定义的production的配置，然后更改
	encoderConfig := zap.NewProductionEncoderConfig()
	//默认时间格式是这样的: "ts":1670214777.9225469 | EpochTimeEncoder serializes a time.Time to a floating-point number of seconds
	//重新设置时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//重新设置时间字段的key
	encoderConfig.TimeKey = "time"
	//默认的level是小写的zapcore.LowercaseLevelEncoder ｜ "level":"info" 可以改成大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//"caller":"zap/zap.go:90" 也可以改成Full的更加详细
	//encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func InitLogger() (err error) {
	//创建核心三大件，进行初始化
	//NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler)
	writerSyncer := getLogWriter(filename, max_size, max_backups, max_age)
	encoder := getEncoder()

	//创建核心
	var core zapcore.Core
	//如果是dev模式，同时要在前端打印；如果是其他模式，就只输出到文件
	if mode == "dev" {
		//使用默认的encoder配置就行了
		//NewConsoleEncoder里面实际上就是一个NewJSONEncoder，需要输入配置
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		//Tee方法将全部日志条目复制到两个或多个底层核心中
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writerSyncer, level),                   //写入到文件的核心
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level), //写到前台的核心
		)
	} else {
		core = zapcore.NewCore(encoder, writerSyncer, level)
	}

	//创建logger对象
	//New方法返回logger，非自定义的情况下就是NewProduction, NewDevelopment,NewExample或者config就可以了。
	//zap.AddCaller是个option，会添加上调用者的文件名和行数，到日志里
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	// return logger
	// 如果return了logger就可以USE之前的ginzap.Ginzap和ginzap.RecoveryWithZap这两个中间件。
	return
}

// GinLogger Gin的middleware，用来替换gin.default的logger
func GinLogger(c *gin.Context) {
	logger := zap.L()
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next() //执行中间件后面的方法，剩余部分等待返回时执行

	cost := time.Since(start) //计算收到请求到回复的时间花费
	logger.Info("request coming to "+path,
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.Duration("cost", cost),
	)
}

// GinRecovery 用于替换gin框架的Recovery中间件，因为传入参数，所以再包一层
// stack bool表示是否记录堆栈信息，可以快速检查错误，但是信息会非常大
func GinRecovery(stack bool) gin.HandlerFunc {
	logger := zap.L()
	return func(c *gin.Context) {
		//defer执行，除了异常，处理并恢复异常，记录日志
		defer func() {
			//这个不必须，检查连接是否断开，因为这个不需要堆栈信息(broken pipe或者connection reset by peer)
			//---------开始--------
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				//httputil包预先准备好的DumpRequest方法
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，我们无法向其写入状态
					c.Error(err.(error))
					c.Abort()
					return
				}
				//---------结束--------

				// 是否打印堆栈信息，使用的是debug.Stack()
				if stack {
					logger.Error(
						"[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error(
						"[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 有错误，直接返回给前端错误，前端直接报错
				//c.AbortWithStatus(http.StatusInternalServerError)
				// 另一种，前端不显示错误
				c.String(200, "访问出错了")
			}
		}()
		c.Next() //这样defer才有效果。先next去执行真正的handler部分，返回的时候才退出这个middleware，那时候才defer
	}
}
