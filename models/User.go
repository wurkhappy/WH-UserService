package Models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/wurkhappy/WH-UserService/DB"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type User struct {
	ID                    bson.ObjectId `bson:"_id"`
	FirstName             string
	LastName              string
	Email                 string
	PwHash                []byte
	AvatarURL             string
	PhoneNumber           string
	DateCreated           time.Time
	AccessToken           string
	AccessTokenSecret     string
	AccessTokenExpiration string
	ConfirmationCode      string
	Fingerprint           string
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password string) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) //this is a panic because bcrypt errors on invalid costs
	}
	u.PwHash = hpass
}

func (u *User) SaveUserWithCtx(ctx *DB.Context) {
	coll := ctx.Database.C("users")
	if err := coll.Insert(u); err != nil {
		return
	}
}

func DeleteUser(id string, ctx *DB.Context) {
	err := ctx.Database.C("users").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Print(err)
		return
	}
}

func NewUser() *User {
	return &User{
		DateCreated: time.Now(),
		ID:          bson.NewObjectId(),
	}
}

func FindUserByEmail(ctx *DB.Context, email string) (u *User, err error) {
	err = ctx.Database.C("users").Find(bson.M{"Email": email}).One(&u)
	if err != nil {
		return
	}
	return
}

func FindUserByID(ctx *DB.Context, id string) (u *User, err error) {
	err = ctx.Database.C("users").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&u)
	if err != nil {
		return
	}
	return
}

func (u *User) PasswordIsValid(password string) {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	if err != nil {
		u = nil
	}
	return
}
