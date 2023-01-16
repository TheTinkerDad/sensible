package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"text/template"

	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/sensors"
	"TheTinkerDad/sensible/settings"
	"TheTinkerDad/sensible/web"
	"TheTinkerDad/sensible/web/api"
)

var Server *http.Server

func apiHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Calling %s...\n", r.URL.Path)

	var result interface{} = nil
	if r.URL.Path == "/api/shutdown" {
		result = api.Shutdown(Server)
	}
	if r.URL.Path == "/api/pause-mqtt" {
		result = api.PauseMqtt()
	}
	if r.URL.Path == "/api/resume-mqtt" {
		result = api.ResumeMqtt()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// Currently not used, but could be used to provide a basic
// control UI or status page, etc.
func pageHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Fetching %s...\n", r.URL.Path)

	var p *web.Page = nil
	if r.URL.Path == "/index.html" {
		p = web.WelcomePage()
	} else {
		p = web.ErrorPage()
	}

	templateString, err := web.WebContent.String("index.html")
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("Index").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, p)
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	staticWebContentHandler := http.FileServer(web.WebContent.HTTPBox())
	http.Handle("/static/", staticWebContentHandler)
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", pageHandler)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Println("Listening on port 8080...")
	return srv
}

func main() {

	settings.EnsureOk()
	mqtt.EnsureOk()
	sensors.EnsureOk()
	web.EnsureOk()

	serverWaitGroup := &sync.WaitGroup{}
	serverWaitGroup.Add(1)
	Server = startHTTPServer(serverWaitGroup)
	serverWaitGroup.Wait()
}
