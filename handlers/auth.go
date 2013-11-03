package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
	// "log"
)

func Login(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	user, _ := models.FindUserByEmail(requestData["email"].(string))

	if user == nil {
		return nil, fmt.Errorf("%s", "Account cannot be found"), http.StatusBadRequest
	}

	if !user.PasswordIsValid(requestData["password"].(string)) {
		return nil, fmt.Errorf("%s", "Invalid password"), http.StatusBadRequest
	}

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}
