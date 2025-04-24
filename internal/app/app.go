package app

import (
	cache "github.com/dijer/otus-go-diploma/internal/cache"
	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/dijer/otus-go-diploma/internal/resizer"
	"github.com/dijer/otus-go-diploma/internal/server"
)

type App struct {
	cache  *cache.Cache
	server server.Server
}

func New(config *config.Config) *App {
	imageCache := cache.New(config.Cache.Size)
	resizer := resizer.New(imageCache, config.Cache.Dir)
	appServer := server.New(config.Server, resizer)

	return &App{
		cache:  &imageCache,
		server: appServer,
	}
}

func (a *App) Run() error {
	return a.server.Start()
}
