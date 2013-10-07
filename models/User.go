package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"github.com/wurkhappy/WH-UserService/DB"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	FirstName   string        `json:"firstName"`
	LastName    string        `json:"lastName"`
	Email       string        `json:"email"`
	PwHash      []byte        `json:"-"`
	AvatarURL   string        `json:"avatarURL"`
	PhoneNumber string        `json:"phoneNumber"`
	DateCreated time.Time     `json:"dateCreated"`
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func NewUser() *User {
	return &User{
		DateCreated: time.Now(),
		ID:          bson.NewObjectId(),
	}
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

func (u *User) SaveUserWithCtx(ctx *DB.Context) (err error) {
	coll := ctx.Database.C("users")
	if _, err := coll.UpsertId(u.ID, &u); err != nil {
		return err
	}
	return nil
}

func FindUserByEmail(email string, ctx *DB.Context) (u *User, err error) {
	err = ctx.Database.C("users").Find(bson.M{"email": email}).One(&u)
	if err != nil {
		return
	}
	return u, nil
}

func FindUserByID(id string, ctx *DB.Context) (u *User, err error) {
	err = ctx.Database.C("users").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteUserWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("users").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func FindUsers(ids []string, ctx *DB.Context) []*User {
	adjustedIDs := make([]bson.ObjectId, 0, len(ids))
	for _, id := range ids {
		adjustedIDs = append(adjustedIDs, bson.ObjectIdHex(id))
	}
	var users []*User
	ctx.Database.C("users").Find(bson.M{"_id": bson.M{"$in": adjustedIDs}}).All(&users)

	return users
}

func (u *User) PasswordIsValid(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	if err != nil {
		return false
	}
	return true
}
