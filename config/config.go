package config

import "bot-middleware/internal/pkg/util"

type RabbitMQConfig struct {
	URL string
}

func LoadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URL: util.GodotEnv("RABBIT_URI"),
	}
}
