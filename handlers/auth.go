package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
)

func Login(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	var requestData map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &requestData)
	user, _ := models.FindUserByEmail(requestData["email"].(string), ctx)

	if user == nil {
		http.Error(w, "Account cannot be found", http.StatusBadRequest)
		return
	}

	if !user.PasswordIsValid(requestData["password"].(string)) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(user)
	w.Write(u)

}

func CreateSignature(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	user, _ := models.FindUserByID(id, ctx)

	var reqData map[string]interface{}
	dec := json.NewDecoder(req.Body)
	dec.Decode(&reqData)

	path := reqData["path"].(string)
	str := user.CreateSignature(path)

	w.Write([]byte(`{"signature":"` + str + `"}`))

}

func VerifySignature(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	user, _ := models.FindUserByID(id, ctx)

	var reqData map[string]interface{}
	dec := json.NewDecoder(req.Body)
	dec.Decode(&reqData)

	path := reqData["path"].(string)
	signature := reqData["signature"].(string)

	if !user.VerifySignature(path, signature) {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}
	user.IsVerified = true
	user.SaveUserWithCtx(ctx)
	u, _ := json.Marshal(user)
	w.Write(u)

}
