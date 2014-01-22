package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
	// "log"
)

func Login(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	user, err := models.FindUserByEmail(requestData["email"].(string))

	if user == nil || err != nil {
		return nil, fmt.Errorf("Sorry, we couldn't find an account with that email and password combination"), http.StatusBadRequest
	}

	if _, ok := requestData["password"]; !ok {
		return nil, fmt.Errorf("Oops! Looks like you forget to enter your password."), http.StatusBadRequest
	}

	if !user.PasswordIsValid(requestData["password"].(string)) {
		return nil, fmt.Errorf("Sorry, we couldn't find an account with that email and password combination"), http.StatusBadRequest
	}

	return user.ToJSON(), nil, http.StatusOK
}
