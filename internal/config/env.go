package config

import (
	_ "embed"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

func NewConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()

	if err != nil {
		slog.Error("viper load config:", err)
	}

	slog.Info("viper setup successfully")
}

func Get(envname string) any {
	return viper.Get(envname)
}

func GetSring(envname string) string {
	return viper.GetString(envname)
}

func GetBool(envname string) bool {
	return viper.GetBool(envname)
}

func GetInt(envname string) int64 {
	return viper.GetInt64(envname)
}

func GetFloat(envname string) float64 {
	return viper.GetFloat64(envname)
}
