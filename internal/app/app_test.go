package app

import (
	"errors"
	"testing"

	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/dijer/otus-go-diploma/internal/server"
	"github.com/stretchr/testify/require"
)

func TestNewApp_Success(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Size: 1024,
			Dir:  "tmp",
		},
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	app := New(cfg)

	require.NotNil(t, app)
	require.NotNil(t, app.cache)
	require.NotNil(t, app.server)
}

func TestApp_RunError(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Size: 1024,
			Dir:  "tmp",
		},
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	mockServer := server.NewMockServer()
	mockServer.On("Start").Return(errors.New("server err"))

	app := New(cfg)
	app.server = mockServer

	err := app.Run()

	require.Error(t, err)
	require.Equal(t, "server err", err.Error())
}
