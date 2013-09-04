package UserController

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func CreateUser(w http.ResponseWriter, req *http.Request) {
	session, err := mgo.Dial("localhost:27017")
        if err != nil {
                panic(err)
        }
        defer session.Close()

        c := session.DB("UserDB").C("users")
        err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	               &Person{"Cla", "+55 53 8402 8510"})
        if err != nil {
                panic(err)
        }
}
