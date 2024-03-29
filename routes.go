package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-UserService/handlers"
)

//order matters so most general should go towards the bottom
var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/user",
			Dest: map[string]interface{}{
				"POST": handlers.CreateUser,
			},
		},
		urlrouter.Route{
			PathExp: "/users",
			Dest: map[string]interface{}{
				"GET": handlers.SearchUsers,
			},
		},
		urlrouter.Route{
			PathExp: "/auth/login",
			Dest: map[string]interface{}{
				"POST": handlers.Login,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/verify",
			Dest: map[string]interface{}{
				"POST": handlers.VerifyUser,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/password",
			Dest: map[string]interface{}{
				"PUT": handlers.NewPassword,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/details",
			Dest: map[string]interface{}{
				"GET": handlers.GetUserDetails,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id",
			Dest: map[string]interface{}{
				"PUT":    handlers.UpdateUser,
				"GET":    handlers.DeleteUser,
				"DELETE": handlers.GetUser,
			},
		},
		urlrouter.Route{
			PathExp: "/password/forgot",
			Dest: map[string]interface{}{
				"POST": handlers.ForgotPassword,
			},
		},
	},
}
