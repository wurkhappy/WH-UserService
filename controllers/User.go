package UserController

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
	user := Models.NewUser()

	var requestData map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &requestData)

	user.Email = requestData["Email"].(string)
	user.SetPassword(requestData["Password"].(string))
	user.FirstName = requestData["FirstName"].(string)
	user.LastName = requestData["LastName"].(string)
	if _, ok := requestData["PhoneNumber"]; ok {
		user.PhoneNumber = requestData["PhoneNumber"].(string)
	}
	user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID.Hex() + ".jpg"
	if _, ok := requestData["AvatarData"]; ok {
		go uploadPhoto(user.ID.Hex(), requestData["AvatarData"].(string))
	}

	user.SaveUserWithCtx(ctx)
	u, _ := json.Marshal(user)
	w.Write(u)

	// fmt.Fprint(w, "Created User")

}

func uploadPhoto(filename string, base64string string) {
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
	_, err = http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	Models.DeleteUser(id, ctx)

	fmt.Fprint(w, "Deleted User")

}
