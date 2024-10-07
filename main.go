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
	//Setup db
	db, err := models.Open(models.DefaultPostgresConfig())
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//Setup services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	//Setup middleware

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}


	key := securecookie.GenerateRandomKey(32)
	csrfMw := csrf.Protect(key, csrf.Secure(false))
	//Setup controllers
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
	//Setup router
	r := chi.NewRouter()
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin",usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Use(csrfMw, umw.SetUser)

	//Start the server
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
