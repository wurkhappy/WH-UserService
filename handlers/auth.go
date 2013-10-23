package handlers

import (
	"encoding/json"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"net/http"
	"fmt"
)

func Login(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	var requestData map[string]interface{}
	json.Unmarshal(body, &requestData)
	user, _ := models.FindUserByEmail(requestData["email"].(string), ctx)

	if user == nil {
		return nil, fmt.Errorf("%s", "Account cannot be found"), http.StatusBadRequest
	}

	if !user.PasswordIsValid(requestData["password"].(string)) {
		return nil, fmt.Errorf("%s", "Invalid password"), http.StatusBadRequest
	}

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}

func CreateSignature(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	id := params["id"].(string)
	user, _ := models.FindUserByID(id, ctx)

	var reqData map[string]interface{}
	json.Unmarshal(body, &reqData)

	path := reqData["path"].(string)
	str := user.CreateSignature(path)

	return []byte(`{"signature":"` + str + `"}`), nil, http.StatusOK

}

func VerifySignature(params map[string]interface{}, body []byte, ctx *DB.Context) ([]byte, error, int) {
	id := params["id"].(string)
	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	var reqData map[string]interface{}
	json.Unmarshal(body, &reqData)

	path := reqData["path"].(string)
	signature := reqData["signature"].(string)

	if !user.VerifySignature(path, signature) {
		return nil, fmt.Errorf("%s", "Invalid signature"), http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}
