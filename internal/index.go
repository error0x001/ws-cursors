package internal

import (
	"html/template"
	"log"
	"net/http"
)

func Index(templatePath, address string, writer http.ResponseWriter, reader *http.Request) {
	log.Println(reader.URL)
	if reader.URL.Path != "/" {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}
	if reader.Method != http.MethodGet {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("parsing template %s error: %s", templatePath, err.Error())
	}
	err = tmpl.Execute(writer, address)
	if err != nil {
		log.Fatalf("exec template %s error: %s", templatePath, err.Error())
	}
}
