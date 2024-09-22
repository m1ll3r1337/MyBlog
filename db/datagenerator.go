package db

import (
	"database/sql"
	"math/rand"
	"strconv"
)

func generateSampleData(db *sql.DB, count int) error {
	userIDs, err := insertUsers(db, count)
	if err != nil {
		return err
	}
	postIDs, err := insertPosts(db, userIDs, count)
	if err != nil {
		return err
	}
	err = insertComments(db, postIDs, userIDs, count)
	if err != nil {
		return err
	}
	imageIDs, err := insertImages(db, count)
	if err != nil {
		return err
	}
	err = linkPostsToImages(db, postIDs, imageIDs)
	if err != nil {
		return err
	}

	return nil
}

func insertUsers(db *sql.DB, count int) ([]int, error) {
	var userIDs []int
	for i := 0; i < count; i++ {
		username := "user" + strconv.Itoa(i)
		email := "user" + strconv.Itoa(i) + "@example.com"
		password := "password" + strconv.Itoa(i)

		var id int
		err := db.QueryRow(`INSERT INTO Users (username, email, password) VALUES ($1, $2, $3) RETURNING id`,
			username, email, password).Scan(&id)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

func insertPosts(db *sql.DB, userIDs []int, count int) ([]int, error) {
	var postIDs []int
	for i := 0; i < count; i++ {
		title := "Post Title " + strconv.Itoa(i)
		content := "This is the content of post number " + strconv.Itoa(i)
		userID := userIDs[rand.Intn(len(userIDs))]

		var id int
		err := db.QueryRow(`INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id`,
			title, content, userID).Scan(&id)
		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, id)
	}
	return postIDs, nil
}

func insertComments(db *sql.DB, postIDs, userIDs []int, count int) error {
	for i := 0; i < count; i++ {
		content := "This is a comment number " + strconv.Itoa(i)
		userID := userIDs[rand.Intn(len(userIDs))]
		postID := postIDs[rand.Intn(len(postIDs))]

		_, err := db.Exec(`INSERT INTO Comments (content, user_id, post_id) VALUES ($1, $2, $3)`,
			content, userID, postID)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertImages(db *sql.DB, count int) ([]int, error) {
	var imageIDs []int
	for i := 0; i < count; i++ {
		data := []byte{0xFF, 0xD8, 0xFF, 0xE0, byte(rand.Intn(255))}

		var id int
		err := db.QueryRow(`INSERT INTO Images (data) VALUES ($1) RETURNING id`, data).Scan(&id)
		if err != nil {
			return nil, err
		}
		imageIDs = append(imageIDs, id)
	}
	return imageIDs, nil
}

func linkPostsToImages(db *sql.DB, postIDs, imageIDs []int) error {
	for i := 0; i < len(postIDs); i++ {
		postID := postIDs[i]
		imageID := imageIDs[rand.Intn(len(imageIDs))]

		_, err := db.Exec(`INSERT INTO Posts_Images (post_id, image_id) VALUES ($1, $2)`,
			postID, imageID)
		if err != nil {
			return err
		}
	}
	return nil
}
