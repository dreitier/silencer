package transport

import (
	"fmt"
	"github.com/go-openapi/runtime/client"
	"silencer/config"
)

func GetTransportClient() *client.Runtime {

	cfg := config.GetInstance().Global()
	url := fmt.Sprintf("%s:%d", cfg.AlertmanagerHost(), cfg.AlertmanagerPort())
	return client.New(url, ApiPath, []string{cfg.AlertmanagerScheme()})
}