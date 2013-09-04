package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/controllers"
	"labix.org/v2/mgo"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", req.URL.Path[1:])
}

func main() {
	DB.Session, err = mgo.Dial(DB.Config.DBURL)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/world", hello).Methods("GET")
	r.HandleFunc("/user", dbContextMixIn(UserController.CreateUser)).Methods("POST")
	http.Handle("/", r)

	http.ListenAndServe(":3000", nil)
}

type dbContextMixIn func(http.ResponseWriter, *http.Request, *Context) error

func (h dbContextMixIn) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//create the context
	ctx, err := NewContext(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.Close()

	//run the handler and grab the error, and report it
	err = h(w, req, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


// I need to Dial the DB in main() and create a session.
// That session needs to get passed to the handlers so that they can clone it.
// Once it's cloned then we can do some ops
// then we have to close the session

