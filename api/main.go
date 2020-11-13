package main

import (
	"github.com/xiao0811/xiniu/route"
)

func main() {
	app := route.GetRouter()

	app.Run()
}
