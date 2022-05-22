package alert

import (
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"silencer/transport"
)

func GetAlerts() (*alert.GetAlertsOK, error) {
	transportClient := transport.GetTransportClient()

	alertClient := alert.New(transportClient, transportClient.Formats)

	params := alert.NewGetAlertsParams()

	return alertClient.GetAlerts(params)
}
