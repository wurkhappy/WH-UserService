package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kr/s3"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"log"
	"net/http"
	"time"
)

var WebServerURL string = "http://localhost:4000"

func CreateUser(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	var err error
	user := models.NewUser()

	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	json.Unmarshal(body, &user)

	if user.Email == "" {
		return nil, fmt.Errorf("%s", "Email cannot be blank"), http.StatusBadRequest
	}

	err = user.SyncWithExistingInvitation()
	if err != nil {
		return nil, err, http.StatusConflict
	}

	pw, ok := requestData["password"].(string)
	if !ok {
		return nil, fmt.Errorf("%s", "Password cannot be blank"), http.StatusConflict
	}

	err = user.ValidateNewPassword(pw)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}
	user.SetPassword(pw)

	if _, ok := requestData["avatarData"]; ok {
		user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID + ".jpg"
		go uploadPhoto(user.ID, requestData["avatarData"].(string))
	} else {
		user.AvatarURL = WebServerURL + "/img/default_photo.jpg"
	}

	user.Save()

	user.AddToPaymentProcessor()
	user.SendVerificationEmail()

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
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
	return resp
}

func GetUser(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	id := params["id"].(string)
	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}

func UpdateUser(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	id := params["id"].(string)
	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	json.Unmarshal(body, &user)

	if pw, ok := requestData["currentPassword"]; ok {
		if !user.PasswordIsValid(pw.(string)) {
			return nil, fmt.Errorf("%s", "Invalid password"), http.StatusBadRequest
		}
		user.SetPassword(requestData["newPassword"].(string))
	}
	if _, ok := requestData["avatarData"]; ok {
		user.AvatarURL = "https://s3.amazonaws.com/PegueNumero/" + user.ID + ".jpg"
		go uploadPhoto(user.ID, requestData["avatarData"].(string))
	}

	user.Save()

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}

func DeleteUser(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	id := params["id"].(string)
	err := models.DeleteUserWithID(id)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}

func SearchUsers(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	var users []*models.User

	if emails, ok := params["email"].([]string); ok {
		for _, email := range emails {
			user, _ := models.FindUserByEmail(email)

			if create, ok := params["create"].([]string); ok && create[0] == "true" && user == nil {
				user = models.NewUser()
				user.Email = email
				user.AvatarURL = WebServerURL + "/img/default_photo.jpg"
				user.Save()
			}
			users = append(users, user)
		}

	}

	if userIDs, ok := params["userid"].([]string); ok {
		users = models.FindUsers(userIDs)
	}

	u, _ := json.Marshal(users)
	return u, nil, http.StatusOK
}

func VerifyUser(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	id := params["id"].(string)
	user, _ := models.FindUserByID(id)

	user.IsVerified = true
	user.Save()

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}

func ForgotPassword(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	data := struct {
		Email string `json:"email"`
	}{}
	json.Unmarshal(body, &data)
	if data.Email == "" {
		return nil, fmt.Errorf("%s", "Email cannot be blank"), http.StatusBadRequest
	}

	user, err := models.FindUserByEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("%s", "There was an error searching for that email"), http.StatusBadRequest
	}

	user.SendForgotPasswordEmail()
	return nil, nil, http.StatusOK
}

func NewPassword(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	var data struct {
		ID       string `json:"id"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}
	data.ID = params["id"].(string)

	json.Unmarshal(body, &data)

	user, err := models.FindUserByID(data.ID)
	if err != nil {
		return nil, fmt.Errorf("%s", "There was an error searching for that user"), http.StatusBadRequest
	}

	if data.Password != data.Confirm {
		return nil, fmt.Errorf("%s", "Passwords do not match"), http.StatusBadRequest
	}
	user.SetPassword(data.Password)
	user.Save()

	return nil, nil, http.StatusOK
}
