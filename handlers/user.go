package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kr/s3"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"log"
	"net/http"
	"time"
)

func CreateUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	user := models.NewUser()

	var requestData map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &requestData)
	json.Unmarshal(buf.Bytes(), &user)

	test, _ := models.FindUserByEmail(requestData["email"].(string), ctx)
	if test != nil {
		if len(test.PwHash) > 0 {
			http.Error(w, "email is already registered", http.StatusConflict)
			return
		} else {
			user.ID = test.ID
		}
	}

	user.SetPassword(requestData["password"].(string))

	if _, ok := requestData["avatarData"]; ok {
		user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID.Hex() + ".jpg"
		go uploadPhoto(user.ID.Hex(), requestData["avatarData"].(string))
	}

	user.AddToPaymentProcessor()
	user.SaveUserWithCtx(ctx)

	user.SendVerificationEmail()

	u, _ := json.Marshal(user)
	w.Write(u)
}

func uploadPhoto(filename string, base64string string) (resp *http.Response) {
	inputFmt := base64string[23 : len(base64string)-1]
	photoData, err := base64.StdEncoding.DecodeString(inputFmt)
	keys := s3.Keys{
		AccessKey: "AKIAJZKKHSBMTOCKBVOA",
		SecretKey: "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",
	}
	data := bytes.NewBuffer(photoData)
	r, _ := http.NewRequest("PUT", "https://s3.amazonaws.com/PegueNumero/"+filename+".jpg", data)
	r.ContentLength = int64(data.Len())
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("X-Amz-Acl", "public-read")
	r.Header.Set("Content-Type", "image/jpeg")
	s3.Sign(r, keys)
	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(resp)
	return resp
}

func GetUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	user, _ := models.FindUserByID(id, ctx)

	u, _ := json.Marshal(user)
	w.Write(u)

}

func UpdateUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	user, _ := models.FindUserByID(id, ctx)

	var requestData map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &requestData)
	json.Unmarshal(buf.Bytes(), &user)

	if pw, ok := requestData["currentPassword"]; ok {
		if !user.PasswordIsValid(pw.(string)) {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}
		user.SetPassword(requestData["newPassword"].(string))
	}
	if _, ok := requestData["avatarData"]; ok {
		user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID.Hex() + ".jpg"
		go uploadPhoto(user.ID.Hex(), requestData["avatarData"].(string))
	}

	user.SaveUserWithCtx(ctx)

	u, _ := json.Marshal(user)
	w.Write(u)

}

func DeleteUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	log.Print(id)
	models.DeleteUserWithID(id, ctx)

	fmt.Fprint(w, "Deleted User")

}

func SearchUsers(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	req.ParseForm()
	var users []*models.User

	//
	if emails, ok := req.Form["email"]; ok {
		for _, email := range emails {
			user, _ := models.FindUserByEmail(email, ctx)

			if create, ok := req.Form["create"]; ok && create[0] == "true" && user == nil {
				user = models.NewUser()
				user.Email = email
				user.SaveUserWithCtx(ctx)
			}
			users = append(users, user)
		}

	}

	if userIDs, ok := req.Form["userid"]; ok {
		log.Print(userIDs)
		users = models.FindUsers(userIDs, ctx)
	}

	u, _ := json.Marshal(users)
	w.Write(u)

}
