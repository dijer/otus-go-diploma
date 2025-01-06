package config

import (
	"os"
	"strconv"
)

type Config struct {
	Cache  CacheConfig
	Server ServerConfig
}

type CacheConfig struct {
	Size int
	Dir  string
}

type ServerConfig struct {
	Port int
	Host string
}

func New() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}
	host := os.Getenv("HOST")

	cacheDir := os.Getenv("CACHE_DIR")
	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Host: host,
		},
		Cache: CacheConfig{
			Dir:  cacheDir,
			Size: cacheSize,
		},
	}, nil
}
