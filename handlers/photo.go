package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/xy3/photos"
	"github.com/xy3/photos/schema"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

// PhotosHandler maps routes on the /photos path
func PhotosHandler(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	if user.ID == 0 {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("photo_id") {
			getPhoto(w, r, user)
		} else {
			getPhotos(w, r, user)
		}
	case http.MethodDelete:
		deletePhoto(w, r, user)
	case http.MethodPut:
		uploadPhoto(w, r, user)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// swagger:route GET /photos/download downloadPhoto
func DownloadPhoto(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	if user.ID == 0 {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	if !r.URL.Query().Has("photo_hash") {
		http.Error(w, "photo_hash query parameter is missing", http.StatusBadRequest)
		return
	}

	photoHash := r.URL.Query().Get("photo_hash")

	var extension string
	row := photos.DB.QueryRow("select extension from photos where file_hash = ? and user_id = ?", photoHash, user.ID)
	err := row.Scan(&extension)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	imagePath := path.Join(user.StoragePath, photoHash, extension)
	http.ServeFile(w, r, imagePath)
}

// getPhoto returns information from the database about a single photo
// swagger:route GET /photos getPhotoInfo
func getPhoto(w http.ResponseWriter, r *http.Request, user schema.User) {
	photoId := r.URL.Query()["photo_id"]
	row := photos.DB.QueryRow("select * from photos where id = ? and user_id = ?", photoId, user.ID)
	photo := &schema.Photo{}
	err := row.Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt, &photo.Size, &photo.FileName, &photo.FileHash, &photo.Extension, &photo.UserID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	marshal, err := json.Marshal(photo)
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}

// getPhotos returns paginated information from the database about multiple photos
// swagger:route GET /photos getPhotos
func getPhotos(w http.ResponseWriter, r *http.Request, user schema.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
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
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	marshal, err := json.Marshal(photoResults)
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}

// swagger:route DELETE /photos deletePhoto
func deletePhoto(w http.ResponseWriter, r *http.Request, user schema.User) {
	tx, err := photos.DB.Begin()
	if err != nil {
		http.Error(w, "failed to begin sql transaction", http.StatusInternalServerError)
		return
	}
	if !r.Form.Has("photo_hash") {
		http.Error(w, "photo_hash query parameter is missing", http.StatusBadRequest)
		return
	}
	photoHash := r.Form.Get("photo_hash")

	// hacky, but ran out of time
	_ = os.Remove(path.Join(user.StoragePath, photoHash+".jpg"))
	_ = os.Remove(path.Join(user.StoragePath, photoHash+".png"))
	_ = os.Remove(path.Join(user.StoragePath, photoHash+".webp"))
	_ = os.Remove(path.Join(user.StoragePath, photoHash+".jpeg"))

	_, err = tx.Exec("delete from photos where file_hash = ? and user_id = ?", photoHash, user.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "failed to delete photo from database", http.StatusInternalServerError)
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("deleted successfully"))
}

// swagger:route PUT /photos uploadPhoto
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
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	sum := md5.Sum(fileBytes)
	fileHash := string(sum[:])
	openFile, err := os.OpenFile(path.Join(user.StoragePath, fileHash), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer openFile.Close()

	tx, _ := photos.DB.Begin()
	_, err = tx.Exec("insert into photos (size, file_name, file_hash, extension, user_id) values (?, ?, ?, ?, ?)", handler.Size, handler.Filename, fileHash, path.Ext(handler.Filename), user.ID)
	if err != nil {
		http.Error(w, "failed to add file to database, photo not created", http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	tx.Commit()

	openFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	w.Write([]byte("success"))
	w.WriteHeader(http.StatusOK)
}
