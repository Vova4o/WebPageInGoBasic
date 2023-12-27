package render

import (
	"WebPageInGoBasic/pkg/config"
	"WebPageInGoBasic/pkg/models"
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// AppConfig sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	// get the template cache from the app config
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the template cache from the app config
	tc = app.TemplateCache

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	// render the template

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// create a template cache map

	myCache := map[string]*template.Template{}

	//get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// loop through the pages one-by-one
	for _, page := range pages {
		// extract the file name (like about.page.tmpl) from the full path
		name := filepath.Base(page)

		// parse the page template file in to a template set
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		// add the template set to the template cache, using the name of the page
		// (like about.page.tmpl) as the key
		myCache[name] = ts
	}
	return myCache, nil
}

//EASY WAY TO PLACE PAGES TO THE CACHE
//var tс = make(map[string]*template.Template)
//
//func RenderTemplate(w http.ResponseWriter, t string) {
//	var tmpl *template.Template
//	var err error
//
//	// check to see if we already have the template in our cache
//	_, ok := tс[t]
//	if !ok {
//		// if not, parse the template and add it to the map
//		log.Println("creating template and adding to cache")
//		err = createTemplateCache(t)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//	} else {
//		// we have the template in the cache
//		log.Println("using cached template")
//	}
//
//	tmpl = tс[t]
//
//	err = tmpl.Execute(w, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//}
//
//func createTemplateCache(t string) error {
//	templates := []string{
//		fmt.Sprintf("./templates/%s", t),
//		"./templates/base.layout.tmpl",
//	}
//
//	// parse the template files...
//	tmpl, err := template.ParseFiles(templates...)
//	if err != nil {
//		return err
//	}
//
//	// ...and add them to the cache
//	tс[t] = tmpl
//	return nil
//}
