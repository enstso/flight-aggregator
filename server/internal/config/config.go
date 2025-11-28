package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	SERVER1_URL string
	SERVER2_URL string
	SERVER_PORT string
)

// Load initializes configuration by reading from a .env file and environment variables, setting relevant global variables.
func Load() {
	// Specify the .env file
	viper.SetConfigFile("../.env")

	// Read the configuration file
	_ = viper.ReadInConfig()

	// Also allow system environment variables (optional)
	viper.AutomaticEnv()

	// Load the values
	SERVER_PORT = viper.GetString("SERVER_PORT")
	j1Name := viper.GetString("JSERVER1_NAME")
	j1Port := viper.GetString("JSERVER1_PORT")
	j2Name := viper.GetString("JSERVER2_NAME")
	j2Port := viper.GetString("JSERVER2_PORT")
	if j1Name != "" && j1Port != "" {
		SERVER1_URL = fmt.Sprintf("http://%s:%s/", j1Name, j1Port)
	}
	if j2Name != "" && j2Port != "" {
		SERVER2_URL = fmt.Sprintf("http://%s:%s/", j2Name, j2Port)
	}

	fmt.Printf("SERVER1_URL: '%s'\n", SERVER1_URL)
	fmt.Printf("SERVER2_URL: '%s'\n", SERVER2_URL)
}
