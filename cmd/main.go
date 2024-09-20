package main

import (
	"github.com/genefriendway/onchain-handler/conf"
	app "github.com/genefriendway/onchain-handler/internal"
)

func main() {
	config := conf.GetConfiguration()
	app.RunApp(config)
}
