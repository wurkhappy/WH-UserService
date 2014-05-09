package models

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-UserService/DB"
	"log"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID                  string    `json:"id" bson:"_id"`
	FirstName           string    `json:"firstName,omitempty"`
	FullFirstName       string    `json:"fullFirstName,omitempty"`
	LastName            string    `json:"lastName,omitempty"`
	Email               string    `json:"email"`
	PwHash              []byte    `json:"-"`
	AvatarURL           string    `json:"avatarURL"`
	PhoneNumber         string    `json:"phoneNumber,omitempty"`
	DateCreated         time.Time `json:"dateCreated"`
	IsVerified          bool      `json:"isVerified"`
	IsProcessorVerified bool      `json:"isProcessorVerified"`
	IsRegistered        bool      `json:"isRegistered"`
	DOBMonth            string    `json:"dobMonth"`
	DOBDay              string    `json:"dobDay"`
	DOBYear             string    `json:"dobYear"`
	SSN                 string    `json:"ssnLastFour"`
	StreetAddress       string    `json:"streetAddress"`
	PostalCode          string    `json:"postalCode"`
}

type user struct {
	ID                  string    `json:"id" bson:"_id"`
	FirstName           string    `json:"firstName,omitempty"`
	FullFirstName       string    `json:"fullFirstName,omitempty"`
	LastName            string    `json:"lastName,omitempty"`
	Email               string    `json:"email"`
	PwHash              []byte    `json:"-"`
	AvatarURL           string    `json:"avatarURL"`
	PhoneNumber         string    `json:"phoneNumber,omitempty"`
	DateCreated         time.Time `json:"dateCreated"`
	IsVerified          bool      `json:"isVerified"`
	IsProcessorVerified bool      `json:"isProcessorVerified"`
	IsRegistered        bool      `json:"isRegistered"`
	DOBMonth            string    `json:"dobMonth"`
	DOBDay              string    `json:"dobDay"`
	DOBYear             string    `json:"dobYear"`
	SSN                 string    `json:"ssnLastFour"`
	StreetAddress       string    `json:"streetAddress"`
	PostalCode          string    `json:"postalCode"`
}

type Users []*User

func (u Users) ToJSON() []byte {
	count := len(u)
	var b bytes.Buffer
	b.WriteString(`[`)
	for i, user := range u {
		b.WriteString(user.toJSONString())
		if i < count-1 {
			b.WriteString(`,`)
		}
	}
	b.WriteString(`]`)
	return b.Bytes()
}

func NewUser() *User {
	id, _ := uuid.NewV4()
	return &User{
		ID:          id.String(),
		DateCreated: time.Now(),
	}
}

func (u *User) SetPassword(password string) error {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PwHash = hpass
	return nil
}

func (u *User) Save() (err error) {
	jsonByte, _ := json.Marshal(u)
	r, err := DB.UpsertUser.Query(u.ID, u.PwHash, string(jsonByte))
	defer r.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (u *User) ToJSON() []byte {
	return []byte(u.toJSONString())
}

func (u *User) toJSONString() string {
	return `{` +
		`"id":"` + u.ID + `",` +
		`"firstName":"` + u.FirstName + `",` +
		`"lastName":"` + u.LastName + `",` +
		`"email":"` + u.Email + `",` +
		`"avatarURL":"` + u.AvatarURL + `",` +
		`"phoneNumber":"` + u.PhoneNumber + `",` +
		`"isVerified":` + strconv.FormatBool(u.IsVerified) + `,` +
		`"isProcessorVerified":` + strconv.FormatBool(u.IsProcessorVerified) + `,` +
		`"isRegistered":` + strconv.FormatBool(u.IsRegistered) + `,` +
		`"dateCreated":` + u.DateCreated.Format(`"`+time.RFC3339Nano+`"`) + `}`
}

func FindUserByEmail(email string) (u *User, err error) {
	var s string
	var b []byte
	err = DB.FindUserByEmail.QueryRow(email).Scan(&b, &s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &u)
	u.PwHash = b
	return u, nil
}

func FindUserByID(id string) (u *User, err error) {
	var s string
	var b []byte
	err = DB.FindUserByID.QueryRow(id).Scan(&b, &s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &u)
	u.PwHash = b
	return u, nil
}

func DeleteUserWithID(id string) (err error) {
	_, err = DB.DeleteUser.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func FindUsers(ids []string) Users {
	for _, id := range ids {
		_, err := uuid.ParseHex(id)
		if err != nil {
			return nil
		}
	}
	fmt.Println(strings.Join(ids, ","))
	var users Users
	r, err := DB.FindUsers.Query("{" + strings.Join(ids, ",") + "}")
	if err != nil {
		log.Print(err)
	}
	defer r.Close()

	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			log.Print(err)
		}
		var u *User
		json.Unmarshal([]byte(s), &u)
		users = append(users, u)
	}
	return users
}

func (u *User) PasswordIsValid(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (u *User) ValidateNewPassword(pw string) error {
	if len(pw) < 6 {
		return fmt.Errorf("%s", "Password length too short")
	}
	return nil
}

func (u *User) IsUserRegistered() bool {
	if len(u.PwHash) > 0 {
		return true
	}
	return false
}

func (u *User) SyncWithExistingInvitation() error {
	test, _ := FindUserByEmail(u.Email)
	if test != nil {
		if test.IsUserRegistered() {
			return fmt.Errorf("%s", "Email is already registered")
		} else {
			u.ID = test.ID
			u.IsVerified = test.IsVerified
		}
	}
	return nil
}

func (u *User) SyncWithExistingUser(existingUser *User) error {
	if existingUser != nil {
		if existingUser.IsUserRegistered() {
			return fmt.Errorf("%s", "Email is already registered")
		} else {
			u.ID = existingUser.ID
			u.IsVerified = existingUser.IsVerified
		}
	}
	return nil
}

func (u *User) UnmarshalJSON(bytes []byte) (err error) {
	var usr *user
	err = json.Unmarshal(bytes, &usr)
	if err != nil {
		return err
	}

	u.ID = usr.ID
	u.FirstName = usr.FirstName
	u.FullFirstName = usr.FullFirstName
	u.LastName = usr.LastName
	u.Email = strings.ToLower(usr.Email)
	u.PwHash = usr.PwHash
	u.AvatarURL = usr.AvatarURL
	u.PhoneNumber = usr.PhoneNumber
	u.DateCreated = usr.DateCreated
	u.IsVerified = usr.IsVerified
	u.IsProcessorVerified = usr.IsProcessorVerified
	u.IsRegistered = usr.IsRegistered
	u.DOBMonth = usr.DOBMonth
	u.DOBDay = usr.DOBDay
	u.DOBYear = usr.DOBYear
	u.SSN = usr.SSN
	u.StreetAddress = usr.StreetAddress
	u.PostalCode = usr.PostalCode
	return nil
}
