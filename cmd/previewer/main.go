package main

import (
	"fmt"
	"os"

	"github.com/dijer/otus-go-diploma/internal/app"
	"github.com/dijer/otus-go-diploma/internal/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		fmt.Println(err)

		return
	}

	err = os.MkdirAll(config.Cache.Dir, 0o755)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := app.New(config)

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
