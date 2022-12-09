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

var dbFileName = "db/deeds.db"

func main() {
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/add", handlerAdd)
	http.HandleFunc("/warning", handlerWarning)
	http.HandleFunc("/update", handlerUpdate)
	http.HandleFunc("/delete", handlerDelete)
	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/styles.css")
	})
	http.HandleFunc("/cookies.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "cookies.html")
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/dev/null")
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
	if !CheckAuth(w, r) {
		return
	}
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
	if !CheckAuth(w, r) {
		return
	}
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
	if !CheckAuth(w, r) {
		return
	}
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
	if !CheckAuth(w, r) {
		return
	}
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

func CheckAuth(w http.ResponseWriter, r *http.Request) bool {
	_, session := Session().Get(r)
	//log.Printf("%s: %v: Session().Get(): %v", id, r.RequestURI, session)
	if session == nil {
		if r.URL.Query()["session"] != nil && r.URL.Query()["session"][0] == "start" {
			http.Redirect(w, r, "/cookies.html", http.StatusSeeOther)
			return false
		}
		Session().Start(w)
		//log.Printf("%v: Session().Start(w)", r.RequestURI)
		http.Redirect(w, r, r.RequestURI+"?session=start", http.StatusSeeOther)
		return false
	}
	if session.Data["auth"] == "auth" {
		return true
	}
	//log.Printf("%s: %s: no auth", id, r.RequestURI)
	password := r.URL.Query().Get("pwd")
	//log.Printf("%s: %s: password = %s", id, r.RequestURI, password)
	if password == "Qr0$21" {
		session.Data["auth"] = "auth"
		return true
	}
	Warning(w, r, "Unauthorised")
	return false
}

/*
	password := r.URL.Query().Get("pwd")
	if password == "Qr0$#21" {
		return true
	}
	//w.WriteHeader(http.StatusUnauthorized)
	Warning(w, r, "Unauthorised")
	return false
*/
