// why do route handlers take in a http.ResponseWriter (and not
// *http.ResponseWriter)?
// http.ResponseWriter is an interface which contains
// a pointer to the underlying content and another pointer to the type
// information. The underlying object is an *http.response which
// is alreday a pointer. Pointers to interfaces are an antipattern in Go

// if http.ReponseWriter is a "pointer" to the underlying object, why not use
// a raw http.response object instead of an *http.response?
// the behavior of Go when passing in an interface which contains a non-pointer
// underlying object as an argument is a copy of the object where the interface
// points to the copy

package main

import (
	// "fmt"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

// Must panics on failure instead of returning an error
// We only parseFiles one time instead of on each viewHandler
// or editHandler call
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// MustCompile panics on failure instead of returning an error
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// render a template using the given page and filename
// write the rendered template to the Response
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	// returns slice or nil if not found
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	// title is second subexpression. m[0] is the whole path
	// and m[1] is either view, edit, or save
	return m[2], nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	// function literal (lambda/anonymous function)
	// closure: a function literal that accesses and/or modifies variables
	// from its surrounding block
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// handler called when path starts with /view/
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// handler called when path starts with /edit/
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// handler called when path starts with /edit/
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	// map routes to handlers
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
