package core

import (
	"fmt"
	"time"

	"backend/global"
	"backend/initialize"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.BE_CONFIG.App.Port)
	s := initServer(address, Router)

	global.BE_LOG.Info("server run success on ", zap.String("address", address))

	// 保证文本顺序输出
	time.Sleep(10 * time.Microsecond)

	fmt.Printf(`
		欢迎使用 多模式访问接口的内存键值存储系统
		当前版本:v0.0.1
		默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
		默认前端文件运行地址:121.40.51.139:8889
	`, address)

	global.BE_LOG.Error(s.ListenAndServe().Error())
}

/*
这些配置可以帮助我们优化HTTP服务器的性能和安全性。通过设置超时时间和最大字节数等参数，可以防止一些潜在的安全问题和性能问题。
例如，设置超时时间可以防止客户端故意保持连接而导致的资源浪费，设置最大字节数可以防止客户端发送过大的请求头而导致的资源浪费和安全问题。
*/
func initServer(address string, router *gin.Engine) server {
	// 使用endless库创建一个HTTP服务器，其中address是服务器的监听地址（如:8080），router是HTTP请求路由器。
	s := endless.NewServer(address, router)
	// 设置HTTP请求头的读取超时时间为20秒，如果在20秒内未读取到请求头，则会返回一个超时错误。
	s.ReadHeaderTimeout = 20 * time.Second

	// 设置HTTP响应体的写入超时时间为20秒，如果在20秒内未将响应体写入完成，则会返回一个超时错误。
	s.WriteTimeout = 20 * time.Second

	// 设置HTTP请求头的最大字节数为1MB。如果请求头超过1MB，则会返回一个错误。
	s.MaxHeaderBytes = 1 << 20

	return s
}
