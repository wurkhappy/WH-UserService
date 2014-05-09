package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
)

func CreateUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var err error
	user := models.NewUser()

	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	json.Unmarshal(body, &user)
	fmt.Println(user.Email, user.FirstName, user.LastName)

	if user.Email == "" {
		return nil, fmt.Errorf("%s", "Email cannot be blank"), http.StatusBadRequest
	}

	//Here we check if the email is already registered (TODO: should break that out into it's own method)
	existingUser, _ := models.FindUserByEmail(user.Email)
	if existingUser != nil {
		err = user.SyncWithExistingUser(existingUser)
		if err != nil {
			return nil, fmt.Errorf("Sorry, could not create the account"), http.StatusConflict
		}
	}

	pw, ok := requestData["password"].(string)
	if !ok {
		return nil, fmt.Errorf("%s", "Password cannot be blank"), http.StatusConflict
	}
	fmt.Println(len(pw))

	err = user.ValidateNewPassword(pw)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}
	user.SetPassword(pw)
	user.IsRegistered = true

	user.Save()

	j := user.ToJSON()
	events := Events{&Event{"user.created", j}}
	go events.Publish()

	return j, nil, http.StatusOK
}

func GetUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	return user.ToJSON(), nil, http.StatusOK
}

func GetUserDetails(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}

func UpdateUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
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

	if user.FirstName == "" {
		user.FirstName = user.FullFirstName
	} else if user.FullFirstName == "" {
		user.FullFirstName = user.FirstName
	}

	user.Save()

	j := user.ToJSON()
	events := Events{&Event{"user.updated", j}}
	go events.Publish()

	return j, nil, http.StatusOK
}

func DeleteUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	err := models.DeleteUserWithID(id)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}

func SearchUsers(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var users models.Users

	if emails, ok := params["email"].([]string); ok {
		for _, email := range emails {
			user, _ := models.FindUserByEmail(email)

			if create, ok := params["create"].([]string); ok && create[0] == "true" && user == nil {
				user = models.NewUser()
				user.Email = email
				user.AvatarURL = "https://d3kq8dzp7eezz0.cloudfront.net/img/default_photo.jpg"
				user.Save()
			}
			users = append(users, user)
		}

	}

	if userIDs, ok := params["userid"].([]string); ok {
		users = models.FindUsers(userIDs)
	}

	return users.ToJSON(), nil, http.StatusOK
}

func VerifyUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	user, _ := models.FindUserByID(id)

	user.IsVerified = true
	user.Save()

	return user.ToJSON(), nil, http.StatusOK
}

func ForgotPassword(params map[string]interface{}, body []byte) ([]byte, error, int) {
	data := struct {
		Email string `json:"email"`
	}{}
	json.Unmarshal(body, &data)
	if data.Email == "" {
		return nil, fmt.Errorf("%s", "Email cannot be blank"), http.StatusBadRequest
	}

	user, err := models.FindUserByEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("%s", "We couldn't find that email. If you need help you can reach us at contact@wurkhappy.com"), http.StatusBadRequest
	}

	j := user.ToJSON()
	events := Events{&Event{"user.forgot_password", j}}
	go events.Publish()
	return nil, nil, http.StatusOK
}

func NewPassword(params map[string]interface{}, body []byte) ([]byte, error, int) {
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
