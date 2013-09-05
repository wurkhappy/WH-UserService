package UserController

import (
	"fmt"
        "github.com/wurkhappy/WH-UserService/models"
        "github.com/wurkhappy/WH-UserService/DB"
	"net/http"
)

func CreateUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	coll := ctx.Database.C("users")

	user := Models.NewUser()
        user.FirstName = "Matt"
        user.LastName = "Parker"
        user.Email = "mdparker89@gmail.com"

	if err := coll.Insert(user); err != nil {
		return
	}

	fmt.Fprint(w, "Created User")

}
