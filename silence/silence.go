package silence

import (
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/prometheus/alertmanager/api/v2/models"
	log "github.com/sirupsen/logrus"
	"silencer/transport"
	"time"
)

func GetSilences() (*silence.GetSilencesOK, error) {
	transportClient := transport.GetTransportClient()
	silenceClient := silence.New(transportClient, transportClient.Formats)
	params := silence.NewGetSilencesParams()

	return silenceClient.GetSilences(params)
}

func FindSilences(filters []string, state string, startsAfter time.Time, endsBefore time.Time) ([]*models.GettableSilence, error) {
	transportClient := transport.GetTransportClient()
	silenceClient := silence.New(transportClient, transportClient.Formats)
	params := silence.NewGetSilencesParams()
	params.Filter = filters

	r, err := silenceClient.GetSilences(params)

	if err != nil {
		return nil, err
	}

	var matchingSilences []*models.GettableSilence

	for _, s := range r.Payload {
		if *s.Status.State != state {
			continue
		}

		startAt, err := time.Parse(time.RFC3339, s.StartsAt.String())

		if err != nil {
			log.Error("failed to parse silence start")
			continue
		}

		endsAt, err := time.Parse(time.RFC3339, s.EndsAt.String())

		if err != nil {
			log.Error("failed to parse silence end")
			continue
		}

		b1 :=  startsAfter.Before(startAt)

			b2 :=endsBefore.After(endsAt)

		if b1 || b2 {
			log.Debug("existing silence starts after or ends before desired time")
			continue
		}

		matchingSilences = append(matchingSilences, s)
	}

	return matchingSilences, err
}

func NewSilence(comment string, createdBy string, startsAt time.Time, endsAt time.Time, matchers models.Matchers) (*silence.PostSilencesOK, error) {
	transportClient := transport.GetTransportClient()
	silenceClient := silence.New(transportClient, transportClient.Formats)
	start, err := strfmt.ParseDateTime(startsAt.Format(time.RFC3339))

	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to parse silence start time")
	}

	end, err := strfmt.ParseDateTime(endsAt.Format(time.RFC3339))

	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to parse silence end time")
	}

	postableSilence := models.PostableSilence{
		ID: "",
		Silence: models.Silence{
			Comment:   &comment,
			CreatedBy: &createdBy,
			EndsAt:    &end,
			Matchers:  matchers,
			StartsAt:  &start,
		},
	}

	params := silence.NewPostSilencesParams()
	params.SetSilence(&postableSilence)

	return silenceClient.PostSilences(params)
}

func DeleteSilence(id strfmt.UUID) (*silence.DeleteSilenceOK, error) {
	transportClient := transport.GetTransportClient()
	silenceClient := silence.New(transportClient, transportClient.Formats)

	params := silence.NewDeleteSilenceParams()
	params.SilenceID = id

	return silenceClient.DeleteSilence(params)
}

func FindOverlappingSilences(startsAt time.Time, endsAt time.Time, matchers models.Matchers) ([]*models.GettableSilence, error) {

	var filters []string

	for _, m := range matchers {
		var operator string

		if *m.IsRegex {
			operator = "=~"
		} else {
			operator = "="
		}

		filter := fmt.Sprintf("%s%s%s", *m.Name, operator, *m.Value)
		filters = append(filters, filter)
	}

	return FindSilences(filters, models.SilenceStatusStateActive, startsAt, endsAt)
}
