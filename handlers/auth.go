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
		return nil, fmt.Errorf("%s", "Account cannot be found"), http.StatusBadRequest
	}

	if _, ok := requestData["password"]; !ok {
		return nil, fmt.Errorf("%s", "Please provide a password"), http.StatusBadRequest
	}

	if !user.PasswordIsValid(requestData["password"].(string)) {
		return nil, fmt.Errorf("%s", "Invalid password"), http.StatusBadRequest
	}

	return user.ToJSON(), nil, http.StatusOK
}
