package main

import (
	"backend/core"
	"backend/global"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const AppMode = "debug" // 运行环境，主要有三种：debug、test、release

func main() {
	gin.SetMode(AppMode)

	//配置初始化
	global.BE_VIPER = core.InitializeViper()

	//日志
	global.BE_LOG = core.InitializeZap()
	zap.ReplaceGlobals(global.BE_LOG)

	global.BE_LOG.Info("server run success on ", zap.String("zap_log", "zap_log"))

	//启动服务
	core.RunServer()
}
