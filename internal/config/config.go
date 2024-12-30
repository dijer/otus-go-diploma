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
	Size int64
	Dir  string
}

type ServerConfig struct {
	Port int64
	Host string
}

func NewConfig() (*Config, error) {
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		return nil, err
	}
	host := os.Getenv("HOST")

	cacheDir := os.Getenv("CACHE_DIR")
	cacheSize, err := strconv.ParseInt(os.Getenv("CACHE_SIZE"), 10, 64)
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
