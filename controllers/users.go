package controllers

import (
	"blog/models"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
		SignIn Template
	}
	Server *Server
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err1 := models.NewUser(username, email, password)
	if err1 != nil {
		 return err1
	}
	id, err := u.Server.Store.CreateUser(user)
	if err != nil {
		return err
	}
	user.ID = id
	WriteJSON(w, http.StatusOK, user)
	return nil
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, data)
}
