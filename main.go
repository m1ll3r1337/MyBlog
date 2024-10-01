package main

import (
	"blog/controllers"
	"blog/models"
	"blog/models/templates"
	"blog/views"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
	"net/http"
)

func main() {
	//store, err := db.NewPostgresStore()
	//if err != nil {
	//	panic(err)
	//}
	//store.Init()
	//err = generateSampleData(store.db, 1000)
	//if err != nil {
	//	panic(err)
	//}
	//s := controllers.Server{Store: store}
	r := chi.NewRouter()
	connStr := "user=postgres dbname=postgres password=balls sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService: &userService,
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
	r.Post("/signup", controllers.MakeHTTPHandleFunc(usersC.Create))
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin",controllers.MakeHTTPHandleFunc(usersC.ProcessSignIn))
	r.Get("/users/me", controllers.MakeHTTPHandleFunc(usersC.CurrentUser))

	//r.Get("/users", controllers.MakeHTTPHandleFunc(s.HandleGetUsers))
	//r.Post("/signup", controllers.MakeHTTPHandleFunc(s.handleCreateUser))
	////r.Get("/users/{id}", makeHTTPHandleFunc(s.handleGetUserByID))
	//r.Put("/users/{id}", controllers.MakeHTTPHandleFunc(s.handleUpdateUser))
	//r.Delete("/users/{id}", controllers.MakeHTTPHandleFunc(s.handleDeleteUser))
	//
	//r.Get("/posts", controllers.MakeHTTPHandleFunc(s.HandleGetComments))
	//r.Post("/posts", controllers.MakeHTTPHandleFunc(s.handleCreatePost))
	//r.Put("/posts/{id}", controllers.MakeHTTPHandleFunc(s.handleUpdatePost))
	//r.Delete("/posts/{id}", controllers.MakeHTTPHandleFunc(s.handleDeletePost))
	//
	//r.Get("/comments", controllers.MakeHTTPHandleFunc(s.HandleGetComments))
	//r.Post("/comments", controllers.MakeHTTPHandleFunc(s.handleCreateComment))
	//r.Put("/comments/{id}", controllers.MakeHTTPHandleFunc(s.handleUpdateComment))
	//r.Delete("/comments/{id}", controllers.MakeHTTPHandleFunc(s.handleDeleteComment))

	//r.Post("/upload", s.HandleUploadImage)
	key := securecookie.GenerateRandomKey(32)
	mw := csrf.Protect(key, csrf.Secure(false))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mw(r))
}
