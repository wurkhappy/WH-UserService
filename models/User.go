package Models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	ID                    bson.ObjectId `bson:"_id,omitempty"`
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

//methods
//passwordIsValid
//Full Name
//setAvatar
//accessTokenIsValid
//confirmationCodeIsValid

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password string) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) //this is a panic because bcrypt errors on invalid costs
	}
	u.PwHash = hpass
}

func NewUser() *User {
	return &User{
		DateCreated: time.Now(),
	}
}
