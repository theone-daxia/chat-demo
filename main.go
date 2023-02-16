package main

import (
	"github.com/theone-daxia/chat-demo/config"
	"github.com/theone-daxia/chat-demo/router"
	"github.com/theone-daxia/chat-demo/service"
)

func main() {
	config.Init()
	go service.Manager.Start()
	r := router.NewRouter()
	_ = r.Run(config.HttpPort)
}
