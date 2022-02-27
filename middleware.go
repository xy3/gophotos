package photos

import (
	"net/http"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()

		if ok {
			user, err := AuthenticateUser(email, password)
			if err == nil {
				EmailToUserCache[email] = user
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
