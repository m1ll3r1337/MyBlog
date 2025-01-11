package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Depado/bfchroma/v2"
	"github.com/microcosm-cc/bluemonday"
	bf "github.com/russross/blackfriday/v2"
	"html/template"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Post struct {
	ID          int
	Title       string
	Content     string
	ContentHTML template.HTML
	UserID      int
	Desc 		string
	Tags 		[]string
}

type Image struct {
	PostID   int
	Path     string
	Filename string
}

type PostService struct {
	DB *sql.DB
	//ImagesDir is used to tell PostService where to store and locate images. If not set, it will default to "images"
	ImagesDir string
	Mr        MarkdownReader
	Mw        MarkdownWriter
	//MarkdownDir is used to tell PostService where to store and locate markdown files. If not set,
	//it will default to "markdowns"
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
	query := `INSERT INTO Posts (title, user_id) VALUES ($1, $2) RETURNING id`
	post := Post{
		Title:  title,
		UserID: userID,
	}
	err := ps.DB.QueryRow(query,
		title,
		userID,
	).Scan(&post.ID)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	dir := ps.postMarkdownDir(post.ID)
	strId := strconv.Itoa(post.ID)
	path := filepath.Join(dir, strId+".md")
	err = ps.Mw.Write(path, content)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return &post, nil
}

func (ps *PostService) GetPaginatedPosts(page int) ([]Post, int, error) {
	totalNumberOfPages, err := ps.getTotalNumberOfPages()
	if err != nil {
		return nil, 0, err
	}
	limit := 20
	offset := (page - 1) * 20

	rows, err := ps.DB.Query(`SELECT id, title, user_id, COALESCE(description, 'No description')
		FROM Posts LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, totalNumberOfPages, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.Title, &post.UserID, &post.Desc); err != nil {
			return nil, totalNumberOfPages, err
		}
		tags, err := ps.GetTagsByPostID(post.ID)
		if err != nil {
			return nil, totalNumberOfPages, err
		}
		post.Tags = tags
		posts = append(posts, post)
	}
	return posts, totalNumberOfPages, nil
}

func (ps *PostService) getTotalNumberOfPages() (int, error) {
	var count int
	query := `SELECT COUNT(id) FROM Posts`
	err := ps.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("getTotalNumberOfPosts: %w", err)
	}
	count = int(math.Ceil(float64(count) / 20))
	return count, nil
}

func (ps *PostService) SearchPosts(urlQuery string) ([]int, error) {
	var postIDs []int
	query := `SELECT 
					DISTINCT (p.id)
				FROM 
					posts p
				LEFT JOIN 
					post_tags pt ON p.id = pt.post_id
				LEFT JOIN tags t ON pt.tag_id = t.id	
				WHERE
					similarity(p.title, $1) > 0.3 OR
					similarity(p.description, $1) > 0.3 OR
					(t.name IS NOT NULL AND similarity(t.name, $1) > 0.3)
				ORDER BY 
					p.id DESC;
				`
	rows, err := ps.DB.Query(query, urlQuery)
	if err != nil {
		return nil, fmt.Errorf("searchPosts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("searchPosts: %w", err)
		}
		postIDs = append(postIDs, id)
	}
	return postIDs, nil
}

func (ps *PostService) GetTagsByPostID(postID int) ([]string, error) {
	query := `SELECT t.name AS tag_name
		FROM Tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = $1;`
	var tags []string
	rows, err := ps.DB.Query(query, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("get tags by post id %d: %w", postID, err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		if err = rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("get tags by post id %d: %w", postID, err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (ps *PostService) GetByID(id int) (*Post, error) {
	query := `SELECT id, title, user_id, COALESCE(description, '') FROM Posts WHERE id = $1`
	var post Post
	err := ps.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.UserID, &post.Desc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	query = `SELECT t.name AS tag_name
		FROM Tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = $1;`
	rows, err := ps.DB.Query(query, id)

	var tags []string
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err = rows.Scan(&tag); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return nil, err
		}
		tags = append(tags, tag)
	}
	post.Tags = tags
	markdownPath, err := ps.Markdown(id)
	if err != nil {
		return nil, err
	}
	post.Content, err = ps.getRawPostMarkdown(markdownPath)
	if err != nil {
		return nil, err
	}
	post.ContentHTML, err = ps.getPostMarkdown(markdownPath)
	return &post, nil
}

func (ps *PostService) Update(post *Post) error {
	_, err := ps.DB.Query(`UPDATE posts SET title=$2, description=$3 WHERE id=$1`,
		post.ID, post.Title, post.Desc)
	query := `SELECT t.name AS tag_name, t.id AS tag_id 
		FROM Tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = $1;`

	var existingTags []struct {
		Name string
		ID   int
	}

	rows, err := ps.DB.Query(query, post.ID)
	defer rows.Close()

	for rows.Next() {
		var tagName string
		var tagID int

		if err = rows.Scan(&tagName, &tagID); err != nil {
			return fmt.Errorf("update post: %w", err)
		}
		existingTags = append(existingTags, struct {
			Name string
			ID   int
		}{
			Name: tagName,
			ID:   tagID,
		})
	}

	for _, tag := range existingTags {
		if !slices.Contains(post.Tags, tag.Name) {
			_, err = ps.DB.Query(`DELETE FROM post_tags WHERE post_id = $1 AND tag_id=$2`, post.ID, tag.ID)
			if err != nil {
				return fmt.Errorf("update post tags: %w", err)
			}
		}
	}
	for _, tag := range post.Tags {
		var tagID int
		err = ps.DB.QueryRow(`INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id;`,
			tag).Scan(&tagID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = ps.DB.QueryRow(`SELECT id FROM tags WHERE name=$1`, tag).Scan(&tagID)
			} else {
				return fmt.Errorf("update post: %w", err)
			}
		}
		var exists bool
		err = ps.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM post_tags WHERE post_id = $1 AND tag_id = $2)`, post.ID, tagID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("update post: %w", err)
		}
		if !exists {
			_, err = ps.DB.Query(`INSERT INTO post_tags (tag_id, post_id) VALUES ($1, $2)`, tagID, post.ID)
			if err != nil {
				return fmt.Errorf("update post: %w", err)
			}
		}

	}
	dir := ps.postMarkdownDir(post.ID)
	strId := strconv.Itoa(post.ID)
	path := filepath.Join(dir, strId+".md")
	post.Content = strings.TrimSpace(post.Content)
	err = ps.Mw.Write(path, post.Content)
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

func (ps *PostService) getPostMarkdown(path string) (template.HTML, error) {
	postMarkdown, err := ps.Mr.Read(path)
	if err != nil {
		return "", err
	}

	b := bf.NewHTMLRenderer(bf.HTMLRendererParameters{
		Flags: bf.CommonHTMLFlags | bf.Smartypants,
	})

	r := bfchroma.NewRenderer(bfchroma.Extend(b), bfchroma.Style("dracula"))
	options := bf.WithExtensions(bf.CommonExtensions | bf.HardLineBreak)

	unsafe := bf.Run([]byte(postMarkdown), options ,bf.WithRenderer(r))
	policy := bluemonday.UGCPolicy()
	policy.AllowElements("br", "code", "pre", "blockquote", "img", "sup", "sub", "strong", "em")
	policy.AllowAttrs("class").OnElements("pre", "code")
	policy.AllowAttrs("src", "alt").OnElements("img")

	html := policy.SanitizeBytes(unsafe)

	return template.HTML(html), nil
}

func (ps *PostService) getRawPostMarkdown(path string) (string, error) {
	postMarkdown, err := ps.Mr.Read(path)
	if err != nil {
		return "", err
	}
	return postMarkdown, nil
}
