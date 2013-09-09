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

	user.SetPassword(requestData["Password"].(string))

	if _, ok := requestData["AvatarData"]; ok {
		user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID.Hex() + ".jpg"
		go uploadPhoto(user.ID.Hex(), requestData["AvatarData"].(string))
	}

	user.SaveUserWithCtx(ctx)

	u, _ := json.Marshal(user)
	w.Write(u)
}

func uploadPhoto(filename string, base64string string) (resp *http.Response) {
	photoData, err := base64.StdEncoding.DecodeString(base64string)
	keys := s3.Keys{
		AccessKey: "AKIAJZKKHSBMTOCKBVOA",
		SecretKey: "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",
	}
	data := bytes.NewBuffer(photoData)
	r, _ := http.NewRequest("PUT", "https://s3.amazonaws.com/PegueNumero/"+filename, data)
	r.ContentLength = int64(data.Len())
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("X-Amz-Acl", "public-read")
	r.Header.Set("Content-Type", "image/jpeg")
	s3.Sign(r, keys)
	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
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
	user := models.NewUser()

	var requestData map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &requestData)
	json.Unmarshal(buf.Bytes(), &user)

	user.SetPassword(requestData["Password"].(string))

	if _, ok := requestData["AvatarData"]; ok {
		user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID.Hex() + ".jpg"
		go uploadPhoto(user.ID.Hex(), requestData["AvatarData"].(string))
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
