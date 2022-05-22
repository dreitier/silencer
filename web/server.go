package web

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func StartServer(port int) {
	listenAddr := fmt.Sprintf(":%d", port)

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(Router)

	srv := &http.Server{
		Handler: r,
		Addr:    listenAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Listening on %s", listenAddr)
	log.Fatal(srv.ListenAndServe())
}
