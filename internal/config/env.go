package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("viper load config:", err)
	}
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
