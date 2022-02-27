package photos

import (
	"github.com/xy3/photos/schema"
	"golang.org/x/crypto/bcrypt"
)

var EmailToUserCache = make(map[string]schema.User)

func getUserByEmail(email string) (user schema.User, err error) {
	row := DB.QueryRow("select * from users where email = ?", email)
	err = row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Password, &user.StoragePath)
	return
}

func AuthenticateUser(email, password string) (user schema.User, err error) {
	user, err = getUserByEmail(email)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return
	}
	return
}