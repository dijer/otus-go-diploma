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

func NewApp(config config.Config) *App {
	imageCache := cache.NewCache(int(config.Cache.Size))
	resizer := resizer.NewResizer(imageCache, config.Cache.Dir)
	appServer := server.NewServer(config.Server, resizer)

	return &App{
		cache:  &imageCache,
		server: appServer,
	}
}

func (a *App) Run() error {
	return a.server.Start()
}
