package main

import (
	log "github.com/sirupsen/logrus"
	"silencer/config"
	"silencer/web"
)

func main() {
	logLevel := config.GetInstance().Global().LogLevel()
	log.SetLevel(logLevel)
	web.StartServer(config.GetInstance().Global().HttpPort())
}
