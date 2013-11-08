package handlers

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func init() {
	config.Test()
	models.Setup()
	DB.Name = "testdb"
	DB.Setup()
	DB.CreateStatements()
	rand.Seed(time.Now().Unix())
}

func generateEmail() string {
	number := rand.Int()
	return strconv.Itoa(number) + "@test.com"
}

func Test(t *testing.T) {
	test_CreateUser(t)
	test_GetUser(t)
	test_UpdateUser(t)
	test_DeleteUser(t)
	test_SearchUsers(t)
	test_VerifyUser(t)
	test_ForgotPassword(t)
	test_NewPassword(t)
	test_Login(t)
	DB.DB.Exec("DELETE from wh_user")
	defer DB.DB.Close()
}

func test_CreateUser(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	params := map[string]interface{}{}

	//
	_, err, _ = CreateUser(params, []byte(""))
	if err == nil {
		t.Error("missing email didn't return error")
	}

	//
	bodyData := map[string]interface{}{
		"email": generateEmail(),
	}
	body, _ := json.Marshal(bodyData)
	_, err, statusCode = CreateUser(params, body)
	if err == nil {
		t.Error("missing password didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "short",
	}
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = CreateUser(params, body)
	if err == nil {
		t.Error("short password didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "short",
	}
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = CreateUser(params, body)
	if err == nil {
		t.Error("short password didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, err, statusCode = CreateUser(params, body)
	if err != nil {
		t.Error("error creating user")
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	var user *models.User
	json.Unmarshal(resp, &user)
	if user.ID == "" {
		t.Error("user wasn't given an ID")
	}
	if user.AvatarURL == "" {
		t.Error("url wasn't set")
	}
}

func test_GetUser(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	params := map[string]interface{}{
		"id": "invalidid",
	}
	_, err, statusCode = GetUser(params, []byte(""))
	if err == nil {
		t.Error("invalid id didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	params = map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, _, _ = CreateUser(params, body)
	var user *models.User
	json.Unmarshal(resp, &user)

	params["id"] = user.ID
	resp, err, statusCode = GetUser(params, []byte(""))
	if err != nil {
		t.Error("error finding user")
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}

func test_UpdateUser(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	params := map[string]interface{}{
		"id": "invalidid",
	}
	_, err, statusCode = UpdateUser(params, []byte(""))
	if err == nil {
		t.Error("invalid id didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	params = map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, _, _ = CreateUser(params, body)
	var user *models.User
	json.Unmarshal(resp, &user)

	params["id"] = user.ID
	bodyData = map[string]interface{}{
		"currentPassword": "wrong",
		"newPassword":     "wrong",
	}
	body, _ = json.Marshal(bodyData)
	resp, err, statusCode = UpdateUser(params, body)
	if err == nil {
		t.Error("did not return error with invalid password")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	params["id"] = user.ID
	bodyData = map[string]interface{}{
		"currentPassword": "password",
		"newPassword":     "newpassword",
		"firstName":       "Tester",
	}
	body, _ = json.Marshal(bodyData)
	resp, err, statusCode = UpdateUser(params, body)
	if err != nil {
		t.Error("error updating user")
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	var u2 *models.User
	json.Unmarshal(resp, &u2)
	if u2.FirstName != bodyData["firstName"].(string) {
		t.Error("name wasn't updated")
	}
}

func test_DeleteUser(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	params := map[string]interface{}{
		"id": "invalidid",
	}
	_, err, statusCode = DeleteUser(params, []byte(""))
	if err == nil {
		t.Error("invalid id didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	params = map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, _, _ = CreateUser(params, body)
	var user *models.User
	json.Unmarshal(resp, &user)

	params["id"] = user.ID
	resp, err, statusCode = DeleteUser(params, []byte(""))
	if err != nil {
		t.Error("error deleting user")
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}

func test_SearchUsers(t *testing.T) {
	// var err error
	var statusCode int
	var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	params := map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, _, _ = CreateUser(params, body)
	var user *models.User
	json.Unmarshal(resp, &user)

	params["email"] = []string{user.Email}
	resp, _, statusCode = SearchUsers(params, []byte(""))
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	var users []*models.User
	json.Unmarshal(resp, &users)
	if len(users) != 1 {
		t.Error("wrong number of users returned")
	}
	if len(users) == 1 && users[0].Email != user.Email {
		t.Error("wrong user returned")
	}

	params["userid"] = []string{user.ID}
	resp, _, statusCode = SearchUsers(params, []byte(""))
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	var users2 []*models.User
	json.Unmarshal(resp, &users2)
	if len(users2) != 1 {
		t.Error("wrong number of users returned")
	}
	if len(users2) == 1 && users2[0].ID != user.ID {
		t.Error("wrong user returned")
	}

	params = map[string]interface{}{}
	params["email"] = []string{generateEmail(), generateEmail(), generateEmail()}
	params["create"] = []string{"true"}
	resp, _, statusCode = SearchUsers(params, []byte(""))
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	var users3 []*models.User
	json.Unmarshal(resp, &users3)
	if len(users3) != 3 {
		t.Error("wrong number of users returned")
	}
	if len(users3) > 0 && users3[0].ID == "" {
		t.Error("created user wasn't given an id")
	}

}

func test_VerifyUser(t *testing.T) {
	var statusCode int
	var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	params := map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, _, _ = CreateUser(params, body)
	var user *models.User
	json.Unmarshal(resp, &user)

	params["id"] = user.ID
	resp, _, statusCode = VerifyUser(params, []byte(""))
	var user2 *models.User
	json.Unmarshal(resp, &user2)
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	if !user2.IsVerified {
		t.Error("user wasn't verified")
	}
}

func test_ForgotPassword(t *testing.T) {
	var err error
	var statusCode int
	var bodyData map[string]interface{}
	var body []byte

	//
	params := map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email": "",
	}
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = ForgotPassword(params, body)
	if err == nil {
		t.Error("blank email didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	bodyData = map[string]interface{}{
		"email": generateEmail(),
	}
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = ForgotPassword(params, body)
	if err == nil {
		t.Error("non-user email didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	params = map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	_, _, _ = CreateUser(params, body)

	_, err, statusCode = ForgotPassword(params, body)
	if err != nil {
		t.Error("error finding user")
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}

func test_NewPassword(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	//
	params := map[string]interface{}{
		"id": "invalidid",
	}
	_, err, statusCode = NewPassword(params, []byte(""))
	if err == nil {
		t.Error("invalid id didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	params = map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": "password",
	}
	body, _ = json.Marshal(bodyData)
	resp, _, _ = CreateUser(params, body)
	var user *models.User
	json.Unmarshal(resp, &user)

	params["id"] = user.ID
	bodyData["confirm"] = "otherpassword"
	body, _ = json.Marshal(bodyData)
	_, _, statusCode = NewPassword(params, body)
	if err == nil {
		t.Error("conflicting passwords didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	params["id"] = user.ID
	bodyData["confirm"] = bodyData["password"].(string)
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = NewPassword(params, body)
	if err != nil {
		t.Errorf("valid params returned error:  %s", err.Error())
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}

func test_Login(t *testing.T) {
	var err error
	var statusCode int
	// var resp []byte
	var bodyData map[string]interface{}
	var body []byte

	//
	//
	params := map[string]interface{}{}
	bodyData = map[string]interface{}{
		"email": generateEmail(),
	}
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = Login(params, body)
	if err == nil {
		t.Error("invalid id didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	//
	params = map[string]interface{}{}
	originalPassword := "password"
	bodyData = map[string]interface{}{
		"email":    generateEmail(),
		"password": originalPassword,
	}
	body, _ = json.Marshal(bodyData)
	CreateUser(params, body)

	bodyData["password"] = "otherpassword"
	body, _ = json.Marshal(bodyData)
	_, _, statusCode = Login(params, body)
	if err == nil {
		t.Error("bad password didn't return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	bodyData["password"] = originalPassword
	body, _ = json.Marshal(bodyData)
	_, err, statusCode = Login(params, body)
	if err != nil {
		t.Errorf("correct password returned error: %s", err.Error())
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}
