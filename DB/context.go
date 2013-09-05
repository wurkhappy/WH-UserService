package DB

import (
	"labix.org/v2/mgo"
	"net/http"
)

type Context struct {
	Database *mgo.Database
}

var Session *mgo.Session

func (c *Context) Close() {
	c.Database.Session.Close()
}

func NewContext(req *http.Request) (*Context, error) {
	return &Context{
		Database: Session.Clone().DB(Config["DBName"]),
	}, nil
}
