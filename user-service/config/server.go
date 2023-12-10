package config

import (
	"fmt"

	"github.com/iqbaludinm/hr-microservice/user-service/utils"
)

type serverConfig struct {
	URI  string
	Port string
	Host string
}

func NewServerConfig() serverConfig {
	return serverConfig{
		URI:  utils.GetEnv("SERVER_URI"),
		Port: utils.GetEnv("SERVER_PORT"),
		Host: fmt.Sprintf("%s:%s", utils.GetEnv("SERVER_URI"), utils.GetEnv("SERVER_PORT")),
	}
}
