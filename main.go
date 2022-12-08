package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kennygrant/sanitize"
)

var dbFileName = "deeds.db"

func main() {
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/add", handlerAdd)
	http.HandleFunc("/warning", handlerWarning)
	http.HandleFunc("/update", handlerUpdate)
	http.HandleFunc("/delete", handlerDelete)
	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/styles.css")
	})
	http.HandleFunc("/fontsizes.css", handlerFontSizes)
	log.Fatal(http.ListenAndServe(":80", nil))
}

type DeedModel struct {
	ID      string
	Name    string
	Color   string
	Overdue string
}

type DeedsModel []*DeedModel

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	db, err := NewDB(dbFileName)
	if err != nil {
		Warning(w, r, "Create DB: %v", err)
		return
	}
	defer db.Close()
	var model DeedsModel
	err = db.Iterate(func(deed *Deed) {
		model = append(model, deed.GetModel())
	})
	if err != nil {
		log.Printf("handler template.New: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templates, err := template.New("page").ParseFiles("web/page.gohtml")
	if err != nil {
		log.Printf("handler template.New: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "page", &model)
	if err != nil {
		log.Printf("handler templates.ExecuteTemplate: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlerAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		Warning(w, r, "%v", err)
		return
	}
	nameStr := r.PostForm.Get("name")
	nameStr = sanitize.HTML(nameStr)
	if nameStr == "" {
		Warning(w, r, "Missing name")
		return
	}
	periodStr := r.PostForm.Get("period")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		Warning(w, r, "Wrong period: \"%s\"", periodStr)
		return
	}
	deed := NewDeed(nameStr, period)
	db, err := NewDB(dbFileName)
	if err != nil {
		Warning(w, r, "Create DB: %v", err)
		return
	}
	defer db.Close()
	err = db.AddDeed(deed)
	if err != nil {
		Warning(w, r, "AddDeed: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if !govalidator.IsUUID(id) {
		Warning(w, r, "Wrong UUID format")
		return
	}
	db, err := NewDB(dbFileName)
	if err != nil {
		Warning(w, r, "Create DB: %v", err)
		return
	}
	defer db.Close()
	err = db.Update(id)
	if err != nil {
		Warning(w, r, "Update: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func handlerDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if !govalidator.IsUUID(id) {
		Warning(w, r, "Wrong UUID format")
		return
	}
	db, err := NewDB(dbFileName)
	if err != nil {
		Warning(w, r, "Create DB: %v", err)
		return
	}
	defer db.Close()
	if err := db.Delete(id); err != nil {
		Warning(w, r, "Delete: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerFontSizes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	f1 := 14
	f2 := 64
	//s1 := 140
	//s2 := 640
	for f := f1; f <= f2; f++ {
		s := f * 10
		//s := s1 + (s2-s1)*(f-f1)/(f2-f1)
		fmt.Fprintf(w, `@media screen and (min-width: %dpx) { 
	* {
		font-size: %dpx;
	}
}

`, s, f)
	}
}
