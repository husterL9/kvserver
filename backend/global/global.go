package global

import (
	"backend/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	BE_CONFIG config.Configuration
	BE_LOG    *zap.Logger
	BE_VIPER  *viper.Viper
)
