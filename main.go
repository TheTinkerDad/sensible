package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/settings"
	"TheTinkerDad/sensible/web"
	"TheTinkerDad/sensible/web/api"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(fmt.Sprintf("Calling %s...", r.URL.Path))

	var object interface{} = nil

	if r.URL.Path == "/api/ansible/info" {
		object = api.AnsibleInfo()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(object)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {

	var p *web.Page = nil

	log.Println(fmt.Sprintf("Fetching %s...", r.URL.Path))

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
	web.EnsureOk()

	serverWaitGroup := &sync.WaitGroup{}
	serverWaitGroup.Add(1)
	srv := startHTTPServer(serverWaitGroup)

	time.Sleep(10 * time.Second)

	log.Println("Shutting down...")
	srv.Shutdown(context.TODO())
	serverWaitGroup.Wait()
}
