package api

import (
	"TheTinkerDad/sensible/settings"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

var Server *http.Server

func apiHandler(w http.ResponseWriter, r *http.Request) {

	var result interface{} = nil

	result = checkToken(r, result)

	log.Debugf("Calling %s...", r.URL.Path)

	if result == nil {
		if r.URL.Path == "/api/info" {
			result = DoInfo()
		} else if r.URL.Path == "/api/shutdown" {
			result = DoShutdown(Server)
		} else if r.URL.Path == "/api/pause-mqtt" {
			result = DoPauseMqtt()
		} else if r.URL.Path == "/api/resume-mqtt" {
			result = DoResumeMqtt()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.(SimpleApiResult).Status)
	json.NewEncoder(w).Encode(result.(SimpleApiResult).Result)
}

func checkToken(r *http.Request, result interface{}) interface{} {
	if settings.All.Api.Token != "" {
		if r.URL.Query().Has("token") {
			if r.URL.Query().Get("token") != settings.All.Api.Token {
				result = ErrWrongToken()
			}
		} else {
			result = ErrMissingToken()
		}
	}
	return result
}

func StartApiServer(wg *sync.WaitGroup) {

	srv := &http.Server{Addr: fmt.Sprintf(":%d", settings.All.Api.Port)}
	http.HandleFunc("/api/", apiHandler)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Infof("Listening on port %d...", settings.All.Api.Port)

	Server = srv
}
