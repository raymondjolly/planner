package render

import (
	"bytes"
	"fmt"
	"github.com/raymondjolly/bookings/internal/config"
	"github.com/raymondjolly/bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

//NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	td.CSRFToken = nosurf.Token(r)
	return td
}

//RenderTemplate renders templates using the html/template package
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		// get the template cache from app config
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		errCheck(err)
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)
	_, err = buf.WriteTo(w)
	errCheck(err)

}

//CreateTemplateCache creates a template cache as a map. It returns either a
//map of *template.Template or an error
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*page.gohtml")
	renderCheck(myCache, err)

	for _, page := range pages {
		name := filepath.Base(page)

		//ts is short for template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		renderCheck(myCache, err)

		matches, err := filepath.Glob("./templates/*layout.gohtml")
		renderCheck(myCache, err)

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.gohtml")
			renderCheck(myCache, err)
		}
		myCache[name] = ts
	}
	return myCache, nil
}

func renderCheck(items map[string]*template.Template, err error) (map[string]*template.Template, error) {
	if err != nil {
		return nil, err
	}
	return items, nil
}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
