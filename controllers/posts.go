package controllers

import (
	"blog/context"
	"blog/models"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

type Posts struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
	}
	PostService *models.PostService
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
	type Image struct {
		PostID          int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID      int
		Title   string
		Content string
		Images  []Image
	}
	data.ID = post.ID
	data.Title = post.Title
	data.Content = post.Content
	images, err := p.PostService.Images(post.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			PostID:          post.ID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
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
		ID              int
		Title           string
		Filename        string
		FilenameEscaped string // TODO: unclear its an image
	}

	var data struct {
		Posts      []Post
	}

	posts, err := p.PostService.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, post := range posts {
		images, err := p.PostService.Images(post.ID) //TODO: multiple image retrieve
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		if len(images) == 0 {
			images = append(images, models.Image{
				PostID:          post.ID,
				Path: "",
				Filename: "",
			})
		}
		data.Posts = append(data.Posts, Post{
			ID:    post.ID,
			Title: post.Title,
			Filename: images[0].Filename,
			FilenameEscaped: url.PathEscape(images[0].Filename),
		})
	}
	p.Templates.Index.Execute(w, r, data)
}

func (p Posts) Show(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r)
	if err != nil {
		return
	}

	type Image struct {
		PostID          int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID      int
		Title   string
		Content string
		Images  []Image
	}
	data.ID = post.ID
	data.Title = post.Title
	data.Content = post.Content
	images, err := p.PostService.Images(post.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			PostID:          post.ID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
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

func (p Posts) Image(w http.ResponseWriter, r *http.Request) {
	filename := p.filename(w, r)
	postID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	image, err := p.PostService.Image(postID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "image not found", http.StatusNotFound)
		}
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, image.Path)
}

func (p Posts) UploadImage(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r, userMustOwnPost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	err = r.ParseMultipartForm(5 << 20) // 5 mb
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = p.PostService.CreateImage(post.ID, fileHeader.Filename, file)
		if err != nil {
			var fileErr models.FileError
			if errors.As(err, &fileErr) {
				msg := fmt.Sprintf("%v has invalid content type or extension. Only png, gif, and jpg files can"+
					"be uploaded.", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
			}
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
	editPath := fmt.Sprintf("/posts/%d/edit", post.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (p Posts) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := p.filename(w, r)
	post, err := p.postByID(w, r, userMustOwnPost)
	if err != nil {
		return
	}
	err = p.PostService.DeleteImage(post.ID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "image not found", http.StatusNotFound)
		}
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts/%d/edit", post.ID), http.StatusFound)
}

func (p Posts) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}

type postOption func(http.ResponseWriter, *http.Request, *models.Post) error

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
