package main

import (
	"blog/controllers"
	"blog/models"
	"blog/models/templates"
	"blog/views"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	db, err := models.Open(models.DefaultPostgresConfig())
	defer db.Close()
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	userService.CreateUsersTable()
	usersC := controllers.Users{
		UserService: &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin",usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	key := securecookie.GenerateRandomKey(32)
	mw := csrf.Protect(key, csrf.Secure(false))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mw(r))
}
