package handlers

import (
	"encoding/json"
	"github.com/xy3/photos"
	"github.com/xy3/photos/schema"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Signup POST /user/signup
func Signup(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorMsg(w, "failed to encrypt password, user not created", http.StatusInternalServerError)
		return
	}
	storagePath, err := photos.CreateUserStorageDirectory(email)
	if err != nil {
		errorMsg(w, "failed to create user storage directory, err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := photos.DB.Begin()
	if err != nil {
		errorMsg(w, "failed to begin sql transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("insert into users (email, password, storage_path) values (?, ?, ?)", email, encryptedPassword, storagePath)
	if err != nil {
		errorMsg(w, "failed to create user, err: "+err.Error(), http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}
	tx.Commit()
	w.WriteHeader(200)
}

// Signin POST /user/signin
func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		invalidMethodError(w)
		return
	}

	_, err := photos.AuthenticateUser(r.PostFormValue("email"), r.PostFormValue("password"))
	if err != nil {
		errorMsg(w, "failed to authenticate user: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// User handles requests on the /user path
func User(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	switch r.Method {
	case http.MethodGet:
		getUser(w, r, user)
	case http.MethodDelete:
		deleteUser(w, r, user)
	default:
		invalidMethodError(w)
	}
}

// getUser GET /user
func getUser(w http.ResponseWriter, _ *http.Request, user schema.User) {
	user.Password = ""
	userJson, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

// deleteUser DELETE /user
func deleteUser(w http.ResponseWriter, _ *http.Request, user schema.User) {
	tx, err := photos.DB.Begin()
	if err != nil {
		errorMsg(w, "failed to begin sql transaction", http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec("delete from users where id = ?", user.ID)
	if err != nil {
		_ = tx.Rollback()
		errorMsg(w, "failed to delete user, err: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(200)
}
