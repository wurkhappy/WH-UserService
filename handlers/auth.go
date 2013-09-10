package handlers

import (
	"crypto/rand"
	"encoding/json"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
	"bytes"
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
	var requestData map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &requestData)
	user, err := models.FindUserByEmail(requestData["Email"].(string), ctx)

	if !user.PasswordIsValid(requestData["Password"].(string)) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(user)
	w.Write(u)

}
