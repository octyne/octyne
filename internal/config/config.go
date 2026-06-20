package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	port := os.Getenv("PORT")
	return Config{
		Port: port,
	}
}
