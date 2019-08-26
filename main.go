package main

import (
	"log"

	"github.com/royge/frameup/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("frameup")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return
	}

	cmd.Execute()
}
