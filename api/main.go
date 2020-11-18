package main

import (
	"log"

	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/route"
)

func main() {
	app := route.GetRouter()

	if err := app.Run(":" + config.Conf.AppConfig.Port); err != nil {
		log.Fatalln(err)
	}
}
