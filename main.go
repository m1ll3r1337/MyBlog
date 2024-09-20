package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	store Storage
}

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		panic(err)
	}
	store.Init()
	//err = generateSampleData(store.db, 1000)
	//if err != nil {
	//	panic(err)
	//}
	s := Server{store: store}
	r := chi.NewRouter()

	r.Get("/users", makeHTTPHandleFunc(s.HandleGetUsers))
	r.Post("/users", makeHTTPHandleFunc(s.handleCreateUser))
	//r.Get("/users/{id}", makeHTTPHandleFunc(s.handleGetUserByID))
	r.Put("/users/{id}", makeHTTPHandleFunc(s.handleUpdateUser))
	r.Delete("/users/{id}", makeHTTPHandleFunc(s.handleDeleteUser))

	r.Get("/posts", makeHTTPHandleFunc(s.HandleGetComments))
	r.Post("/posts", makeHTTPHandleFunc(s.handleCreatePost))
	r.Put("/posts/{id}", makeHTTPHandleFunc(s.handleUpdatePost))
	r.Delete("/posts/{id}", makeHTTPHandleFunc(s.handleDeletePost))

	r.Get("/comments", makeHTTPHandleFunc(s.HandleGetComments))
	r.Post("/comments", makeHTTPHandleFunc(s.handleCreateComment))
	r.Put("/comments/{id}", makeHTTPHandleFunc(s.handleUpdateComment))
	r.Delete("/comments/{id}", makeHTTPHandleFunc(s.handleDeleteComment))

	r.Post("/upload", s.HandleUploadImage)
	http.ListenAndServe(":8080", r)
	fmt.Println("Listening on port 8080")
}
