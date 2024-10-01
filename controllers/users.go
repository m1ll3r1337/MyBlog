package controllers

import (
	"blog/models"
	"fmt"
	"net/http"
	"strings"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	email = strings.ToLower(email)
	user, err1 := models.NewUser(username, email, password)
	if err1 != nil {
		return err1
	}
	id, err := u.UserService.CreateUser(user)
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
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) error {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	WriteJSON(w, http.StatusOK, user)

	return nil
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("email")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, cookie.Value)
	return nil
}
