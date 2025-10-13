package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	SERVER1_URL string
	SERVER2_URL string
)

// Load initializes configuration by reading from a .env file and environment variables, setting relevant global variables.
func Load() {
	// Specify the .env file
	viper.SetConfigFile(".env")

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading .env file: %v\n", err)
		return
	}

	// Also allow system environment variables (optional)
	viper.AutomaticEnv()

	// Load the values
	SERVER1_URL = viper.GetString("SERVER1_URL")
	SERVER2_URL = viper.GetString("SERVER2_URL")

	fmt.Printf("SERVER1_URL: '%s'\n", SERVER1_URL)
	fmt.Printf("SERVER2_URL: '%s'\n", SERVER2_URL)
}
