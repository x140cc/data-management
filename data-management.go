package main

import (
	"flag"
	"github.com/scompo/data-management/domain"
	"html/template"
	"log"
	"net/http"
)

const appName = "data-management"

type Page struct {
	Title    string
	PageName string
}

var port = flag.String("port", "8080", "server port")

func main() {

	flag.Parse()

	log.Printf("Initializing %v...\n", appName)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/projects/new", newProjectHandler)

	log.Printf("listening on port %v...", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func newProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called newProjectHandler")
	switch r.Method {
	case "POST":
		log.Printf("method: POST\n")
		r.ParseForm()
		name := r.FormValue("Name")
		description := r.FormValue("Description")
		pi := domain.ProjectInfo{
			Project: domain.Project{
				Name: name,
			},
			Description: description,
		}
		log.Printf("project %v\n", pi)
		domain.SaveProject(pi)
		http.Redirect(w, r, "/projects", http.StatusFound)
		return
	case "GET":
		log.Printf("method: GET\n")
		t, err := template.ParseFiles(
			"templates/main.html",
			"templates/header.html",
			"templates/new-project.html")
		p := Page{Title: appName, PageName: "projects"}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, map[string]interface{}{
			"Page": p,
		})
		return
	default:
		log.Printf("DEFAULT\n")
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called projectsHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/projects.html")
	p := Page{Title: appName, PageName: "projects"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"Page":     p,
		"Projects": domain.AllProjects(),
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called mainHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/index.html")
	p := Page{Title: appName, PageName: "index"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"Page": p,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
