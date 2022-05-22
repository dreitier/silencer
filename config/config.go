package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

type configuration struct {
	global       *global
}

var (
	instance *configuration
	once     sync.Once
	cfgPaths []string
)

const (
	CfgFileName = "config.yaml"
	PathLocal   = "."
	PathGlobal  = "/etc/silencer"

)

func init() {

	cfgPaths = append(cfgPaths, PathLocal)

	userHome, err := os.UserHomeDir()

	if  err == nil{
		userHome = fmt.Sprintf("%s%c%s", userHome, os.PathSeparator, ".silencer")
		cfgPaths = append(cfgPaths, userHome)
	}

	cfgPaths = append(cfgPaths, PathGlobal)
}

func GetInstance() *configuration {
	once.Do(func() {
		instance = &configuration{}
		initConfig()
	})
	return instance
}

func (c *configuration) Global() *global {
	return c.global
}

func initConfig() {

	var file *os.File = nil
	var err error = nil

	for _, path := range cfgPaths {
		file, err = os.Open(filepath.Join(path, CfgFileName))

		if err == nil {
			log.Infof("found config file at location %s", path)
			break
		}
	}

	if file == nil {
		log.Fatal("could not open config file")
	}

	defer file.Close()

	cfg, err := Parse(file)
	if err != nil {
		log.Fatalf("failed to parse config file: %s", err)
	}

	instance.global = parseGlobal(cfg)
}

func parseGlobal(cfg Raw) *global {

	g := newGlobal()

	if cfg.Has("log_level") {
		g.logLevel = cfg.String("log_level")
	}

	if cfg.Has("port") {
		g.httpPort = int(cfg.Int64("port"))
	}

	if cfg.Has("alertmanager_host") {
		g.alertmanagerHost = cfg.String("alertmanager_host")
	}

	if cfg.Has("alertmanager_scheme") {
		g.alertmanagerScheme = cfg.String("alertmanager_scheme")
	}

	if cfg.Has("alertmanager_port") {
		g.alertmanagerPort = int(cfg.Int64("alertmanager_port"))
	}
	
	if cfg.Has("known_services") {
		g.knownServices = cfg.StringSlice("known_services")
	}

	return g
}


