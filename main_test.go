package main

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	// "github.com/wurkhappy/WH-UserService/models"
	"bytes"
	"encoding/json"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/controllers"
	"io"
	"labix.org/v2/mgo"

)

type nopCloser struct {
	io.Reader
}

func (n nopCloser) Close() error {
	return nil
}



func TestCreateUser(t *testing.T) {
	//test id is returned
	//test that first name is the same
	//test that pass
	// DB.Session, _ = mgo.Dial(DB.Config["DBURL"])

	// user := map[string]interface{}{
	// 	"FirstName": "Test",
	// 	"Password":  "password",
	// }
	// u, _ := json.Marshal(user)

	// record := httptest.NewRecorder()
	// req := &http.Request{
	// 	Method: "POST",
	// 	URL:    &url.URL{Path: "/user"},
	// 	Body:   nopCloser{bytes.NewBuffer(u)},
	// }

	// //create the context
	// ctx, _ := DB.NewContext(req)
	// defer ctx.Close()

	// Controllers.CreateUser(record, req, ctx)
	// // if got, want := record.Code, 400; got != want {
	// // 	t.Errorf("%s: response code = %d, want %d", "Hello world test", got, want)
	// // }
	// // fmt.Print("HI")
}
