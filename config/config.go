package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

// AppConf main application configuration
type AppConf struct {
	// Listener configuration
	Listener struct {
		Host string `yaml:"host" env:"HOST" env-default:"localhost"` // Host to be listened
		Port int    `yaml:"port" env:"PORT" env-default:"7373"`      // Port to be listened
	} `yaml:"listener"`

	SecretKey string `yaml:"secret_key" env:"SECRET_KEY" env-default:"super_secret_key!"`

	// Target destination hosts
	Destinations []string `yaml:"destinations"`
}

// Prepare listen address string
func (c *AppConf) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Listener.Host, c.Listener.Port)
}

// Init for application config initialization
// If filePath provided, then default values will be overwritten by values in file and after that by environment variables
// If filePath not provided default values will be overwritten by environment variables only
func Init(filePath string) *AppConf {
	var conf AppConf

	var err error
	if filePath != "" {
		log.Info().Msgf("Reading config from file %s", filePath)
		err = cleanenv.ReadConfig(filePath, &conf)
	} else {
		log.Info().Msgf("Reading config from environment variables")
		err = cleanenv.ReadEnv(&conf)
	}

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot parse config")
	}

	return &conf
}
