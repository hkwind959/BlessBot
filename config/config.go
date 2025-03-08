package config

import (
	"BlessBot/logs"
	"BlessBot/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var fileConfig *viper.Viper

func Init() {
	logs.Log().Info("初始化配置文件...")
	viper.SetConfigFile("./config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		logs.Log().Fatal("读取配置文件失败", zap.Error(err))
	}
	fileConfig = viper.GetViper()
}

// GetConfig 获取配置文件
func GetConfig() model.Account {
	var AppConfig model.Account
	err := fileConfig.Unmarshal(&AppConfig)
	if err != nil {
		logs.Log().Fatal("解析配置文件失败", zap.Error(err))
	}
	logs.Log().Info("配置文件获取成功")
	return AppConfig
}
