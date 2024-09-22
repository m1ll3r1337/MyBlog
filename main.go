package main

import (
	"blog/controllers"
	"blog/db"
	"blog/templates"
	"blog/views"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	store, err := db.NewPostgresStore()
	if err != nil {
		panic(err)
	}
	store.Init()
	//err = generateSampleData(store.db, 1000)
	//if err != nil {
	//	panic(err)
	//}
	s := controllers.Server{Store: store}
	r := chi.NewRouter()

	usersC := controllers.Users{
		Server: &s,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))

	r.Get("/signup", usersC.New)
	r.Post("/signup", controllers.MakeHTTPHandleFunc(s.HandleCreateUser))

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
	http.ListenAndServe(":8080", r)
	fmt.Println("Listening on port 8080")
}
