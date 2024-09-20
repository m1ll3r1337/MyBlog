package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func (s *Server) HandleGetUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *Server) HandleGetPosts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	posts, err := s.store.GetPosts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, posts)
}

func (s *Server) HandleGetComments(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	comments, err := s.store.GetComments()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, comments)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := new(CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}
	if !isValidEmail(createUserReq.Email) {
		_ = fmt.Errorf("invalid email")
	}
	if len(createUserReq.Password) < 8 {
		err := fmt.Errorf("password needs to be atleast 8 characters long")
		return err
	}
	if len(createUserReq.Username) == 0 {
		err := fmt.Errorf("empty username")
		return err
	}

	user, err := NewUser(createUserReq.Username, createUserReq.Email, createUserReq.Password)
	if err != nil {
		return err
	}

	id, err := s.store.CreateUser(user)
	if err != nil {
		return err
	}
	user.ID = id

	return WriteJSON(w, http.StatusOK, user)
}

func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	createPostReq := new(CreatePostRequest)

	if err := json.NewDecoder(r.Body).Decode(createPostReq); err != nil {
		return err
	}

	if len(createPostReq.Title) == 0 {
		err := fmt.Errorf("empty title")
		return err
	}
	if len(createPostReq.Content) == 0 {
		err := fmt.Errorf("empty content")
		return err
	}

	post := NewPost(createPostReq.Title, createPostReq.Content, createPostReq.UserID)

	id, err := s.store.CreatePost(post)
	if err != nil {
		return err
	}
	post.ID = id

	return WriteJSON(w, http.StatusOK, post)
}

func (s *Server) handleCreateComment(w http.ResponseWriter, r *http.Request) error {
	createCommentReq := new(CreateCommentRequest)

	if err := json.NewDecoder(r.Body).Decode(createCommentReq); err != nil {
		return err
	}

	if len(createCommentReq.Content) == 0 {
		err := fmt.Errorf("empty content")
		return err
	}

	comment := NewComment(createCommentReq.Content, createCommentReq.UserID, createCommentReq.PostID)

	id, err := s.store.CreateComment(comment)
	if err != nil {
		return err
	}
	comment.ID = id

	return WriteJSON(w, http.StatusOK, comment)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	id, err1 := getID(r)
	if err1 != nil {
		return err1
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return err
	}
	if !isValidEmail(user.Email) {
		err = fmt.Errorf("invalid email")
		return err
	}
	if len(user.Username) == 0 {
		err = fmt.Errorf("empty username")
		return err
	}

	if err = s.store.UpdateUser(id, &user); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": id})
}

func (s *Server) handleUpdatePost(w http.ResponseWriter, r *http.Request) error {
	id, err1 := getID(r)
	if err1 != nil {
		return err1
	}
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return err
	}

	if len(post.Title) == 0 {
		err = fmt.Errorf("empty title")
		return err
	}
	if len(post.Content) == 0 {
		err = fmt.Errorf("empty content")
		return err
	}

	if err = s.store.UpdatePost(id, &post); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": id})
}

func (s *Server) handleUpdateComment(w http.ResponseWriter, r *http.Request) error {
	id, err1 := getID(r)
	if err1 != nil {
		return err1
	}
	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		return err
	}

	if len(comment.Content) == 0 {
		err = fmt.Errorf("empty content")
		return err
	}


	if err = s.store.UpdateComment(id, &comment); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": id})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err = s.store.DeleteUser(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *Server) handleDeletePost(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err = s.store.DeletePost(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *Server) handleDeleteComment(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err = s.store.DeleteComment(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *Server) HandleUploadImage(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	err = s.store.UploadImage(fileBytes)
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		fmt.Println("Error matching regular expression:", err)
		return false
	}
	return matched
}

func getID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}