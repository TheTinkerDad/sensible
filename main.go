package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
        "os/exec"
	"strconv"
	"text/template"

	"github.com/GeertJohan/go.rice"
	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Title string
	Body  string
}

func savePlaybooks() {

	database, err := sql.Open("sqlite3", "./settings.db")
        if err != nil {
                log.Fatal(err)
        }

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS playbooks (id INTEGER PRIMARY KEY, filename TEXT, description TEXT)")
        if err != nil {
                log.Fatal(err)
        }
	statement.Exec()

	statement, err = database.Prepare("INSERT INTO playbooks (filename, description) VALUES (?, ?)")
        if err != nil {
                log.Fatal(err)
        }
	statement.Exec("shutdown-node.yaml", "Restarts one or more nodes")

	rows, err := database.Query("SELECT id, filename, description FROM playbooks")
        if err != nil {
                log.Fatal(err)
        }

	var id int
	var filename string
	var desc string
	for rows.Next() {
		rows.Scan(&id, &filename, &desc)
		log.Println(strconv.Itoa(id) + ": " + filename + " " + desc)
	}
}

var box *rice.Box = nil

func loadPage() *Page {
	return &Page{Title: "Hello", Body: "Hello World!"}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {

	var p *Page = nil

	log.Println(fmt.Sprintf("Fetching %s...", r.URL.Path))

	if r.URL.Path == "/index.html" {
		p = loadPage()
	} else if r.URL.Path == "/info" {
                out, _ := exec.Command("ansible", "--version").Output()
		p = &Page{Title: "Ansible Info", Body: string(out)}
	}

	templateString, err := box.String("index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("Index").Parse(templateString)
        if err != nil {
                log.Fatal(err)
        }

	t.Execute(w, p)
}

func main() {

	savePlaybooks()

	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}
	box, _ = conf.FindBox("data")

	cssFileServer := http.FileServer(box.HTTPBox())
	http.Handle("/css/", cssFileServer)
	http.HandleFunc("/", pageHandler)
        log.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
