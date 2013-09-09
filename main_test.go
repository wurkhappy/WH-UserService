package main

import (
	// "fmt"
	"bytes"
	"encoding/json"
	// "github.com/gorilla/context"
	// "github.com/gorilla/mux"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/controllers"
	"github.com/wurkhappy/WH-UserService/models"
	"io"
	"labix.org/v2/mgo"
	// "log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type nopCloser struct {
	io.Reader
}

func (n nopCloser) Close() error {
	return nil
}

func ClearDB() {

}

func NewContext(req *http.Request) (*DB.Context, error) {
	return &DB.Context{
		Database: DB.Session.Clone().DB("TestUserDB"),
	}, nil
}

func TestCreateUser(t *testing.T) {
	DB.Session, _ = mgo.Dial(DB.Config["DBURL"])

	userParams := map[string]interface{}{
		"FirstName": "Test",
		"Password":  "password",
	}
	u, _ := json.Marshal(userParams)

	record := httptest.NewRecorder()
	req := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/user"},
		Body:   nopCloser{bytes.NewBuffer(u)},
	}

	ctx, _ := NewContext(req)
	defer ctx.Close()

	Controllers.CreateUser(record, req, ctx)
	if gotCode, wantCode := record.Code, 200; gotCode != wantCode {
		t.Errorf("%s:%d RESULT: response code = %d", "Should respond with code", wantCode, gotCode)
	}

	user := new(models.User)
	decoder := json.NewDecoder(record.Body)
	decoder.Decode(&user)
	if gotName, wantName := user.FirstName, userParams["FirstName"]; gotName != wantName {
		t.Errorf("%s:%d RESULT: FirstName = %d", "Should return first name", wantName, gotName)
	}

	if !user.ID.Valid() {
		t.Error("Did not return a user id")
	}
}
