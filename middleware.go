package photos

import (
	"github.com/xy3/photos/schema"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var EmailToUserCache = make(map[string]schema.User)

func BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()
		if ok {
			row := DB.QueryRow("select * from users where email = ?", email)
			var user schema.User

			err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.StoragePath)
			if err == nil {
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err == nil {
					EmailToUserCache[email] = user
					next.ServeHTTP(w, r)
					return
				}
			}

		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
