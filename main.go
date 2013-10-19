package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/handlers"
	"labix.org/v2/mgo"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", req.URL.Path[1:])
}

func main() {
	var err error
	DB.Session, err = mgo.Dial(DB.Config["DBURL"])
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/world", hello).Methods("GET")
	r.Handle("/user", dbContextMixIn(handlers.CreateUser)).Methods("POST")
	r.Handle("/user/search", dbContextMixIn(handlers.SearchUsers)).Methods("GET")
	r.Handle("/auth/login", dbContextMixIn(handlers.Login)).Methods("POST")

	//these two don't feel RESTful. Will have to think about it some more
	r.Handle("/user/{id}/sign", dbContextMixIn(handlers.CreateSignature)).Methods("POST")
	r.Handle("/user/{id}/sign/verify", dbContextMixIn(handlers.VerifySignature)).Methods("POST")

	r.Handle("/user/{id}", dbContextMixIn(handlers.UpdateUser)).Methods("PUT")
	r.Handle("/user/{id}/verify", dbContextMixIn(handlers.VerifyUser)).Methods("POST")
	r.Handle("/user/{id}", dbContextMixIn(handlers.DeleteUser)).Methods("DELETE")
	r.Handle("/user/{id}", dbContextMixIn(handlers.GetUser)).Methods("GET")
	http.Handle("/", r)

	http.ListenAndServe(":3000", nil)
}

type dbContextMixIn func(http.ResponseWriter, *http.Request, *DB.Context)

func (h dbContextMixIn) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//create the context
	ctx, err := DB.NewContext(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.Close()

	h(w, req, ctx)
}
