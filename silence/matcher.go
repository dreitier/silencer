package silence

import "github.com/prometheus/alertmanager/api/v2/models"

func NewMatcher(isRegex bool, name string, value string) *models.Matcher {
	return &models.Matcher{
		IsRegex: &isRegex,
		Name:    &name,
		Value:   &value,
	}
}
