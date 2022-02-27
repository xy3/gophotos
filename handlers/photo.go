package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/xy3/photos"
	"github.com/xy3/photos/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	DefaultPageSize = 10
	DefaultPage     = 0
)

// Photo handles requests on the /photo path
func Photo(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	switch r.Method {
	case http.MethodGet:
		getPhoto(w, r, user)
	case http.MethodDelete:
		deletePhoto(w, r, user)
	case http.MethodPut:
		uploadPhoto(w, r, user)
	default:
		invalidMethodError(w)
	}
}

// PhotoInfo handles requests on the /photo/info path
func PhotoInfo(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	switch r.Method {
	case http.MethodGet:
		getPhotoInfo(w, r, user)
	case http.MethodPatch:
		updatePhotoInfo(w, r, user)
	default:
		invalidMethodError(w)
	}
}

// PhotoList handles requests on the /photo/list path
func PhotoList(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	switch r.Method {
	case http.MethodGet:
		getPhotoList(w, r, user)
	default:
		invalidMethodError(w)
	}
}

// updatePhotoInfo PATCH /photo/info
// Update information about a photo on the server
func updatePhotoInfo(w http.ResponseWriter, r *http.Request, user schema.User) {
	photo := schema.Photo{}
	err := json.NewDecoder(r.Body).Decode(&photo)
	if err != nil {
		errorMsg(w, "could not decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if photo.ID == 0 {
		errorMsg(w, "photo ID is missing", http.StatusBadRequest)
		return
	}

	var exists bool
	result := photos.DB.QueryRow("select exists (select * from photos where id = ? and user_id = ?)", photo.ID, user.ID)
	result.Scan(&exists)
	if !exists {
		errorMsg(w, "the photo id provided does not match a photo in your account", http.StatusNotFound)
		return
	}

	tx, err := photos.DB.Begin()
	if err != nil {
		errorMsg(w, "could not initialize database transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("update photos set size = ?, file_name = ?, file_hash = ?, extension = ?, updated_at = CURRENT_TIMESTAMP where id = ? and user_id = ?", photo.Size, photo.FileName, photo.FileHash, photo.Extension, photo.ID, user.ID)
	if err != nil {
		_ = tx.Rollback()
		errorMsg(w, "failed to update photo information: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tx.Commit()
	w.WriteHeader(http.StatusOK)
}

// getPhoto GET /photo
func getPhoto(w http.ResponseWriter, r *http.Request, user schema.User) {
	if !r.URL.Query().Has("photo_id") {
		errorMsg(w, "photo_hash query parameter is missing", http.StatusBadRequest)
		return
	}

	photoId := r.URL.Query().Get("photo_id")

	var (
		fileHash  string
		extension string
	)

	row := photos.DB.QueryRow("select file_hash, extension from photos where id = ? and user_id = ?", photoId, user.ID)
	err := row.Scan(&fileHash, &extension)
	if err != nil {
		errorMsg(w, "not found", http.StatusNotFound)
		return
	}

	imagePath := path.Join(user.StoragePath, fileHash+extension)
	http.ServeFile(w, r, imagePath)
}

// getPhotoInfo GET /photo/info
// Returns information from the database about a single photo
func getPhotoInfo(w http.ResponseWriter, r *http.Request, user schema.User) {
	if !r.URL.Query().Has("photo_id") {
		errorMsg(w, "missing photo_id parameter", http.StatusBadRequest)
		return
	}
	photoId := r.URL.Query().Get("photo_id")
	row := photos.DB.QueryRow("select * from photos where id = ? and user_id = ?", photoId, user.ID)
	photo := schema.Photo{}
	err := row.Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt, &photo.Size, &photo.FileName, &photo.FileHash, &photo.Extension, &photo.UserID)
	if err != nil {
		errorMsg(w, "not found: "+err.Error(), http.StatusNotFound)
		return
	}
	photoJson, _ := json.Marshal(photo)
	w.WriteHeader(http.StatusOK)
	w.Write(photoJson)
}

// getPhotoList GET /photo getPhotoList
// Returns paginated information from the database about multiple photos
func getPhotoList(w http.ResponseWriter, r *http.Request, user schema.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = DefaultPage
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = DefaultPageSize
	}

	photoResults := make([]*schema.Photo, 0)
	rows, err := photos.DB.Query("select * from photos where user_id = ? limit ? offset ?", user.ID, pageSize, page*pageSize)
	defer rows.Close()
	for rows.Next() {
		photo := new(schema.Photo)
		_ = rows.Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt, &photo.Size, &photo.FileName, &photo.FileHash, &photo.Extension, &photo.UserID)
		photoResults = append(photoResults, photo)
	}

	if err = rows.Err(); err != nil {
		errorMsg(w, "not found", http.StatusNotFound)
		return
	}
	photosJson, err := json.Marshal(photoResults)
	if err != nil {
		errorMsg(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(photosJson)
}

// deletePhoto DELETE /photo
func deletePhoto(w http.ResponseWriter, r *http.Request, user schema.User) {
	if !r.URL.Query().Has("photo_id") {
		errorMsg(w, "missing photo_id parameter", http.StatusBadRequest)
		return
	}
	photoId := r.URL.Query().Get("photo_id")

	row := photos.DB.QueryRow("select file_hash, extension from photos where id =?", photoId)
	var (
		fileHash  string
		extension string
	)

	err := row.Scan(&fileHash, &extension)
	if err != nil {
		errorMsg(w, "failed to find photo file in the database", http.StatusNotFound)
		return
	}

	err = os.Remove(path.Join(user.StoragePath, fileHash+extension))
	if err != nil {
		errorMsg(w, "failed to delete the photo file on the disk: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := photos.DB.Begin()
	if err != nil {
		errorMsg(w, "failed to begin sql transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("delete from photos where id = ? and user_id = ?", photoId, user.ID)
	if err != nil {
		_ = tx.Rollback()
		errorMsg(w, "failed to delete photo from database", http.StatusInternalServerError)
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

// uploadPhoto PUT /photo
// Upload a photo to the server
func uploadPhoto(w http.ResponseWriter, r *http.Request, user schema.User) {
	// 10 << 20 = max 10MB upload
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("photo")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	sum := md5.Sum(fileBytes)
	fileHash := fmt.Sprintf("%x", sum[:])
	openFile, err := os.OpenFile(path.Join(user.StoragePath, fileHash+path.Ext(handler.Filename)), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer openFile.Close()

	tx, _ := photos.DB.Begin()
	_, err = tx.Exec("insert into photos (size, file_name, file_hash, extension, user_id) values (?, ?, ?, ?, ?)", handler.Size, handler.Filename, fileHash, path.Ext(handler.Filename), user.ID)
	if err != nil {
		errorMsg(w, "failed to add file to database, photo not created: "+err.Error(), http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}

	if _, err = openFile.Write(fileBytes); err != nil {
		errorMsg(w, "failed to write file data to disk: "+err.Error(), http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}

	newPhoto := schema.Photo{}
	result := tx.QueryRow("select * from photos where file_hash = ?", fileHash)
	result.Scan(&newPhoto.ID, &newPhoto.CreatedAt, &newPhoto.UpdatedAt, &newPhoto.Size, &newPhoto.FileName, &newPhoto.FileHash, &newPhoto.Extension, &newPhoto.UserID)

	_ = tx.Commit()
	fmt.Printf("Successfully uploaded File: %+v\n", newPhoto)

	photoJson, _ := json.Marshal(newPhoto)
	w.Write(photoJson)
}
