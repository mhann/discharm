package cmd

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	log.Println("Initializing viper configuration")
	
	viper.SetConfigName("config") // name of the config file (without extension)
	viper.AddConfigPath("/etc/discharm/")
	viper.AddConfigPath("$HOME/.discharm")
	viper.AddConfigPath(".")
	
	log.Println("Initializing default configuration values")
	
	viper.SetDefault("BotID", "Bot MjczMTI0NjMxNjU2MDA1NjMz.C26-Ug.tZrN1HhotClAem-yQTlNsleKFbE")
	viper.SetDefault("TwitchChannelChecks", map[string]string{})
	
	log.Println("Reading in configuration file")
	
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}
