package controllers

import (
	"blog/context"
	"blog/models"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type Posts struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
	}
	PostService  *models.PostService
	ImageService *models.ImageService
}

func (p Posts) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title   string
		Content string
	}
	data.Title = r.FormValue("title")
	data.Content = r.PostFormValue("content")
	p.Templates.New.Execute(w, r, data)
}

func (p Posts) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID  int
		Title   string
		Content string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	data.Content = r.FormValue("content")

	post, err := p.PostService.Create(data.Title, data.Content, data.UserID)
	if err != nil {
		p.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("posts/%d/edit", post.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (p Posts) Edit(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r, userMustOwnPost)
	if err != nil {
		return
	}

	data := struct {
		ID      int
		Title   string
		Content string
	}{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	}
	p.Templates.Edit.Execute(w, r, data)
}

func (p Posts) Update(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r, userMustOwnPost)
	if err != nil {
		return
	}

	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	err = p.PostService.Update(post)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts/%d/edit", post.ID), http.StatusFound)
}

func (p Posts) Index(w http.ResponseWriter, r *http.Request) {
	type Post struct {
		ID    int
		Title string
	}
	var data struct {
		Posts []Post
	}

	posts, err := p.PostService.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, post := range posts {
		data.Posts = append(data.Posts, Post{
			ID:    post.ID,
			Title: post.Title,
		})
	}
	p.Templates.Index.Execute(w, r, data)
}

func (p Posts) Show(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r)
	if err != nil {
		return
	}

	image, err := p.ImageService.ByPostID(post.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	data := struct {
		ID      int
		Title   string
		Content string
		Image   *models.Image
	}{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Image: image,
	}
	p.Templates.Show.Execute(w, r, data)
}

func (p Posts) Delete(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r, userMustOwnPost)
	if err != nil {
		return
	}

	err = p.PostService.Delete(post.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

type postOption func(http.ResponseWriter,*http.Request, *models.Post) error

func (p Posts) postByID(w http.ResponseWriter, r *http.Request, opts ...postOption) (*models.Post, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusNotFound)
		return nil, err
	}
	post, err := p.PostService.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Invalid id", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	for _, opt := range opts {
		err = opt(w, r, post)
		if err != nil {
			return nil, err
		}
	}
	return post, nil
}

func userMustOwnPost(w http.ResponseWriter, r *http.Request, post *models.Post) error {
	user := context.User(r.Context())
	if post.UserID != user.ID {
		http.Error(w, "You are not authorized to edit the post", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}
	return nil
}
