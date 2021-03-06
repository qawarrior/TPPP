package handler

import (
	"html/template"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/qawarrior/serve-nt/model"
)

type profile struct {
	users  *model.UsersCollection
	events *model.EventsCollection
}

func (h *profile) get(w http.ResponseWriter, r *http.Request) {
	if !authenicated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// obtain and validate the id
	id := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(id) == false {
		cfg.Logger.Error.Println("The ID:", id, "is not valid")
		http.Error(w, "Invalid Id", http.StatusUnauthorized)
		return
	}

	oid := bson.ObjectIdHex(id)
	u, err := h.users.FindOne(map[string]interface{}{"_id": oid})
	if err != nil {
		cfg.Logger.Error.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	page := model.PageData{
		Timestamp: time.Now(),
		AppName:   cfg.AppName,
	}

	evts, err := h.events.Find(map[string]interface{}{
		"serveeid": oid,
	})
	if err != nil {
		cfg.Logger.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := model.ProfileData{
		PageData: page,
		User:     *u,
		Events:   evts,
	}
	tpl, err := template.ParseFiles("./assets/templates/_layout.html", "./assets/templates/profile.html")
	if err != nil {
		cfg.Logger.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tpl.ExecuteTemplate(w, "_layout", p)
}
