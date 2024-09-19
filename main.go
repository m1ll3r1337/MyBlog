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
	if err != nil {
		panic(err)
	}
	s := Server{store: store}
	r := chi.NewRouter()
	r.Get("/users", s.HandleGetUsers)
	http.ListenAndServe(":8080", r)
	fmt.Println("Listening on port 8080")
}
