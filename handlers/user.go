package handlers

import (
	"encoding/json"
	"github.com/xy3/photos"
	"github.com/xy3/photos/schema"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type signupRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,required"`
	Password string `json:"password,required"`
}

// swagger:route POST /users/signup signup
func Signup(w http.ResponseWriter, r *http.Request) {
	request := signupRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "signup request is invalid", http.StatusBadRequest)
		return
	}
	tx, err := photos.DB.Begin()
	if err != nil {
		http.Error(w, "failed to begin sql transaction", http.StatusInternalServerError)
		return
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to encrypt password, user not created", http.StatusInternalServerError)
		return
	}
	storagePath, err := photos.CreateUserStorageDirectory(request.Email)
	if err != nil {
		http.Error(w, "failed to create user storage directory, err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("insert into users (name, email, password, storage_path) values (?, ?, ?, ?)", request.Name, request.Email, encryptedPassword, storagePath)
	if err != nil {
		http.Error(w, "failed to create user, err: "+err.Error(), http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	tx.Commit()
	w.WriteHeader(200)
	w.Write([]byte("user created successfully"))
}

// UsersHandler maps routes on the /users path
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	user := photos.EmailToUserCache[email]
	if user.ID == 0 {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getUser(w, r, user)
	case http.MethodDelete:
		deleteUser(w, r, user)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// swagger:route GET /users getUser
func getUser(w http.ResponseWriter, r *http.Request, user schema.User) {
	marshal, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "failed to marshal json data", http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}

// swagger:route DELETE /users deleteUser
func deleteUser(w http.ResponseWriter, _ *http.Request, user schema.User) {
	tx, err := photos.DB.Begin()
	if err != nil {
		http.Error(w, "failed to begin sql transaction", http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec("delete from users where id = ?", user.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "failed to delete user, err: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(200)
	w.Write([]byte("deleted successfully"))
}
