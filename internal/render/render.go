package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/imsujan276/golang-bookings/internal/config"
	"github.com/imsujan276/golang-bookings/internal/models"
	"github.com/justinas/nosurf"
)

//
var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(tempData *models.TemplateData, r *http.Request) *models.TemplateData {
	tempData.Flash = app.Session.PopString(r.Context(), "flash")
	tempData.Warning = app.Session.PopString(r.Context(), "warning")
	tempData.Error = app.Session.PopString(r.Context(), "error")
	tempData.CSRFToken = nosurf.Token(r)
	return tempData
}

// RenderTemplate renders the static templates form templates directory
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tempData *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	temp, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	// take the template, execute it, dont pass data and store it in buf variable
	tempData = AddDefaultData(tempData, r)
	temp.Execute(buf, tempData)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error parseing template:", tmpl)
		return
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	// err = parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	log.Println("Error parseing template:", tmpl)
	// 	return
	// }
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// myCache is a map of string of type pointer of template.Template{}
	myCache := map[string]*template.Template{}

	// search for all the pages inside templates folder
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.html")
		}

		myCache[name] = templateSet
	}
	return myCache, nil
}
