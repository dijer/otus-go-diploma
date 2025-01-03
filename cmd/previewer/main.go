package main

import (
	"fmt"
	"os"

	"github.com/dijer/otus-go-diploma/internal/app"
	"github.com/dijer/otus-go-diploma/internal/config"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll(config.Cache.Dir, 0755); err != nil {
		fmt.Println(err)
		return
	}

	app := app.NewApp(*config)

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
