package config

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"golang.org/x/exp/slices"
)

type global struct {
	logLevel           string
	httpPort           int
	knownServices      []string
	alertmanagerHost   string
	alertmanagerPort   int
	alertmanagerScheme string
}

func newGlobal() *global {
	return &global{
		logLevel:           "info",
		httpPort:           8000,
		knownServices:      []string{"*"},
		alertmanagerHost:   "alertmanager.prometheus.svc",
		alertmanagerPort:   80,
		alertmanagerScheme: "http",
	}
}

func (g *global) LogLevel() log.Level {
	switch strings.ToLower(g.logLevel) {
	case "error":
		return log.ErrorLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	default:
		return log.InfoLevel

	}
}

func (g *global) HttpPort() int {
	return g.httpPort
}

func (g *global) AlertmanagerHost() string {
	return g.alertmanagerHost
}

func (g *global) AlertmanagerPort() int {
	return g.alertmanagerPort
}

func (g *global) AlertmanagerScheme() string {
	return g.alertmanagerScheme
}

func (g *global) IsServiceKnown(service string) bool {
	// accept every service
	if (slices.Contains(g.knownServices, "*")) {
		return true
	}

	return slices.Contains(g.knownServices, service)
}
