package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/kennygrant/sanitize"
)

func Warning(w http.ResponseWriter, r *http.Request, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	url := fmt.Sprintf("/warning?text=%s", message)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

type WarningModel struct {
	Text string
}

func handlerWarning(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	text = sanitize.HTML(text)
	model := WarningModel{Text: text}

	templates, err := template.New("page").ParseFiles("web/main_layout.gohtml", "web/warning_view.gohtml")
	if err != nil {
		log.Printf("handlerWarning template.New: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "page", &model)
	if err != nil {
		log.Printf("handlerWarning templates.ExecuteTemplate: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
