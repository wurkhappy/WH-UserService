package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
	"strings"
	// "log"
)

func Login(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	fmt.Println(requestData["email"])

	email, ok := requestData["email"].(string)

	if !ok {
		return nil, fmt.Errorf("Please include an email address"), http.StatusBadRequest
	}

	user, err := models.FindUserByEmail(strings.ToLower(email))
	if user == nil || err != nil {
		return nil, fmt.Errorf("Sorry, we couldn't find an account with that email and password combination"), http.StatusBadRequest
	}

	password, ok := requestData["password"].(string)

	if !ok {
		return nil, fmt.Errorf("Oops! Looks like you forget to enter your password."), http.StatusBadRequest
	}

	if !user.PasswordIsValid(password) {
		return nil, fmt.Errorf("Sorry, we couldn't find an account with that email and password combination"), http.StatusBadRequest
	}

	return user.ToJSON(), nil, http.StatusOK
}
