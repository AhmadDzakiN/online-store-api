package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func NewViperConfig() *viper.Viper {
	dirPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("error get working dir: %s", err))
	}
	dirPaths := strings.Split(dirPath, "/internal")

	godotenv.Load(fmt.Sprintf("%s/params/.env", dirPaths[0]))
	godotenv.Load("./params/.env")

	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("../../../../params")
	v.AddConfigPath("./params")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err == nil {
		fmt.Printf("Using config file: %s \n", v.ConfigFileUsed())
	} else {

		panic(fmt.Errorf("Config error: %s", err.Error()))
	}

	return v
}
