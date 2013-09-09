package handlers

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/wurkhappy/WH-UserService/models"
	"http"
	"time"
)

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func Login(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	user := models.FindUserByEmail(email, ctx)

	if !user.PasswordIsValid(password) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(user)
	w.Write(u)

}
