package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-UserService/DB"
	// "labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID          string    `json:"id" bson:"_id"`
	FirstName   string    `json:"firstName,omitempty"`
	LastName    string    `json:"lastName,omitempty"`
	Email       string    `json:"email"`
	PwHash      []byte    `json:"-"`
	AvatarURL   string    `json:"avatarURL"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	DateCreated time.Time `json:"dateCreated"`
	IsVerified  bool      `json:"isVerified"`
}

func NewUser() *User {
	id, _ := uuid.NewV4()
	return &User{
		ID: id.String(),
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
	_, err = DB.UpsertUser.Query(u.ID, u.PwHash, string(jsonByte))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
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

func FindUsers(ids []string) []*User {
	for _, id := range ids {
		_, err := uuid.ParseHex(id)
		if err != nil {
			return nil
		}
	}
	var users []*User
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

func (u *User) SaveUserWithCtx(ctx *DB.Context) (err error) {
	coll := ctx.Database.C("users")
	if _, err := coll.UpsertId(u.ID, &u); err != nil {
		return err
	}
	return nil
}

func (u *User) PasswordIsValid(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (u *User) AddToPaymentProcessor() {
	client := &http.Client{}
	r, _ := http.NewRequest("POST", PaymentInfoService+"/user/"+u.ID, nil)
	_, err := client.Do(r)
	if err != nil {
	}
}

func (u *User) SendVerificationEmail() {
	message := map[string]interface{}{
		"Body": map[string]interface{}{
			"user": u,
		},
	}
	body, _ := json.Marshal(message)
	sendEmail("/user/verify", body)
}

func (u *User) SendForgotPasswordEmail() {
	message := map[string]interface{}{
		"Body": map[string]interface{}{
			"user": u,
		},
	}
	body, _ := json.Marshal(message)
	sendEmail("/user/password/forgot", body)
}

func sendEmail(path string, body []byte) error {
	publisher, err := rbtmq.NewPublisher(connection, emailExchange, "direct", emailQueue, path)
	if err != nil {
		return err
	}
	publisher.Publish(body, false)
	return nil
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
		}
	}
	return nil
}
