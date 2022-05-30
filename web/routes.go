package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/alertmanager/api/v2/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"github.com/dreitier/silencer/silence"
	"github.com/dreitier/silencer/config"
	"strconv"
	"time"
)

var Router *mux.Router

func init() {
	Router = mux.NewRouter().UseEncodedPath()
	Router.StrictSlash(true)

	Router.HandleFunc("/silence/{service}/{tenant:\\w+}-{stage:\\w+}", SilenceHandler).Methods("POST")
}

func SilenceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	unescape(vars)
	durationParam := r.FormValue("duration")
	comment := r.FormValue("comment")
	if comment == ""{
		comment = "Created by Silencer"
	}

	if durationParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("Required parameter 'duration' not present")))
		return
	}

	duration, err := strconv.ParseInt(durationParam, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("form parameter duration could not be parsed", err)
		_, _ = w.Write([]byte(fmt.Sprint("Parameter 'duration' could not be parsed")))
		return
	}

	service := vars["service"]
	tenant := vars["tenant"]
	stage := vars["stage"]

	
	if config.GetInstance().Global().IsServiceKnown(service) == false {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Service is not known"))
		return
	}

	startAt := time.Now()
	endAt := startAt.Add(time.Second * time.Duration(duration))


	matchers := models.Matchers{}
	matchers = append(matchers, silence.NewMatcher(false, "tenant", tenant))
	matchers = append(matchers, silence.NewMatcher(false, "stage", stage))


	overlappingSilences, err := silence.FindOverlappingSilences(startAt, endAt, matchers)

	if err != nil {
		log.Error("failed to find overlapping silences, continuing anyway", err)
	}

	if len(overlappingSilences) > 0 {
		log.Infof("found active silence(s) with compatible matchers already covering the requested duration, no new silence is created")
		_, _ = w.Write([]byte(fmt.Sprintf("no new silence is created, because request is already covered by existing silence with id %s", *overlappingSilences[0].ID)))
		return
	}


	silenceOk, err := silence.NewSilence(comment, "Silencer", startAt, endAt, matchers)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprint("Failed to create silence", err)))
		return
	}

	silenceOkBody, err := silenceOk.Payload.MarshalBinary()
	log.Debugf("Silence created with id ", string(silenceOkBody))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprint("Failed to marshal result from Alertmanager", err)))
		return
	}
	response := fmt.Sprintf("Silence created. Tenant: %s, stage: %s, duration: %ds", tenant, stage, duration)
	log.Infof(response)
	_, _ = w.Write([]byte(response))
}

func unescape(vars map[string]string) {
	for key, val := range vars {
		val, err := url.PathUnescape(val)
		if err == nil {
			vars[key] = val
		}
	}
}
