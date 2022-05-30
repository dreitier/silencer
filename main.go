package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/dreitier/silencer/config"
	"github.com/dreitier/silencer/web"
)

func main() {
	logLevel := config.GetInstance().Global().LogLevel()
	log.SetLevel(logLevel)
	web.StartServer(config.GetInstance().Global().HttpPort())
}
