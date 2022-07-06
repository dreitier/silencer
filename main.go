package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/dreitier/silencer/config"
	"github.com/dreitier/silencer/web"
)

const app = "silencer"
var gitRepo = "dreitier/silencer"
var gitCommit = "unknown"
var gitTag = "unknown"

func printVersion() {
	if gitTag == "" {
		gitTag = "err-no-git-tag"
	}

	log.Printf("%s (dist=%s; version=%s; commit=%s)", app, gitRepo, gitTag, gitCommit)
}

func main() {
	printVersion()

	logLevel := config.GetInstance().Global().LogLevel()
	log.SetLevel(logLevel)
	web.StartServer(config.GetInstance().Global().HttpPort())
}
