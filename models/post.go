package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type Image struct {
	PostID   int
	Path     string
	Filename string
}



type PostService struct {
	DB *sql.DB
	//ImagesDir is used to tell PostService where to store and locate images. If not set, it will default to "images"
	ImagesDir   string
	Sr          MarkdownReader
	MarkdownDir string
}

func NewPost(title, content string, userID int) *Post {
	return &Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
}

func (ps *PostService) Create(title, content string, userID int) (*Post, error) {
	query := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id`
	post := Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
	err := ps.DB.QueryRow(query,
		title,
		content,
		userID,
	).Scan(&post.ID)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return &post, nil
}

func (ps *PostService) GetAll() ([]Post, error) {
	rows, err := ps.DB.Query(`SELECT id, title, content, user_id FROM Posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (ps *PostService) GetByID(id int) (*Post, error) {
	query := `SELECT id, title, content, user_id FROM Posts WHERE id = $1`
	var post Post
	err := ps.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (ps *PostService) Update(post *Post) error {
	_, err := ps.DB.Query(`UPDATE posts SET title=$2, content=$3 WHERE id=$1`,
		post.ID, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	return nil
}

func (ps *PostService) Images(postID int) ([]Image, error) {
	globPattern := filepath.Join(ps.postImagesDir(postID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving images: %w", err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, ps.imageExtensions()) {
			images = append(images, Image{
				PostID:   postID,
				Path:     file,
				Filename: filepath.Base(file),
			})
		}
	}
	return images, nil
}

func (ps *PostService) Image(postID int, filename string) (Image, error) {
	imagePath := filepath.Join(ps.postImagesDir(postID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("image not found: %w", err)
	}
	return Image{
		Filename: filename,
		PostID:   postID,
		Path:     imagePath,
	}, nil
}

//func (ps *PostService) Markdown(postID int) (Markdown, error) {
//	markdownPath := filepath.Join(ps.postMarkdownDir(postID))
//	_, err := os.Stat(markdownPath)
//	if err != nil {
//		if errors.Is(err, fs.ErrNotExist) {
//			return Markdown{}, ErrNotFound
//		}
//		return Image{}, fmt.Errorf("image not found: %w", err)
//	}
//	return Image{
//		Filename: filename,
//		PostID:   postID,
//		Path:     imagePath,
//	}, nil
//}

func (ps *PostService) Delete(id int) error {
	_, err := ps.DB.Exec(`delete from posts where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	dir := ps.postImagesDir(id)
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}

func (ps *PostService) CreateImage(postID int, filename string, contents io.ReadSeeker) error {
	// TODO: how to keep a single image at all times
	err := checkContentType(contents, ps.imageContentTypes())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	err = checkExtension(filename, ps.imageExtensions())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	postDir := ps.postImagesDir(postID)
	imagePath := filepath.Join(postDir, filename)
	err = os.MkdirAll(postDir, 0755)
	if err != nil {
		return fmt.Errorf("creating image directory: %w", err)
	}
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("writing image file: %w", err)
	}

	return nil
}

func (ps *PostService) DeleteImage(postID int, filename string) error {
	image, err := ps.Image(postID, filename)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
}

func (ps *PostService) Markdown(postID int) (string, error) {
	markdownDir := ps.postMarkdownDir(postID)
	filename := fmt.Sprintf("%d.md", postID)
	markdownPath := filepath.Join(markdownDir, filename)
	_, err := os.Stat(markdownPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("markdown file not found: %w", err)
	}
	return markdownPath, nil
}

func (ps *PostService) postImagesDir(id int) string {
	imagesDir := ps.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("post-%d", id))
}

func (ps *PostService) postMarkdownDir(id int) string {
	markdownDir := ps.MarkdownDir
	if markdownDir == "" {
		markdownDir = "markdowns"
	}
	return markdownDir
}


func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (ps *PostService) imageExtensions() []string {
	return []string{".png", ".jpg", ".gif", ".jpeg"}
}

func (ps *PostService) imageContentTypes() []string {
	return []string{"image/png", "image/jpg", "image/gif", "image/jpeg"}
}
