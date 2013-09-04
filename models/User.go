package models

import (
	"code.google.com/p/go.crypto/bcrypt"
)

type User struct{
	FirstName string
	LastName string
	Email string
	PwHash string
	AvatarURL string
	PhoneNumber string
	DateCreated string
	AccessToken string
	AccessTokenSecret string
	AccessTokenExpiration string
	ConfirmationCode string
	ID string
	Fingerprint string
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
	u.Password = hpass
}
