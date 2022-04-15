package config

import (
	//"time"
	"time"

	"github.com/spf13/viper"
)

var (
	InitConfig = initConfig
)

func initConfig(fileName string, additionalDirs []string) (err error) {
	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	for _, dir := range additionalDirs {
		viper.AddConfigPath(dir)
	}
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	viper.ConfigFileUsed()
	viper.WatchConfig()
	return
}

func getConfigString(key string) string {
	return viper.GetString(key)
}

func getConfigInt(key string) int {
	return viper.GetInt(key)
}

func getConfigDuration(key string) time.Duration {
	return viper.GetDuration(key)
}