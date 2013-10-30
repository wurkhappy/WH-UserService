package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-UserService/DB"
	"labix.org/v2/mgo/bson"
	// "log"
	"net/http"
	"time"
)

type User struct {
	ID          string    `json:"id" bson:"_id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PwHash      []byte    `json:"-"`
	AvatarURL   string    `json:"avatarURL"`
	PhoneNumber string    `json:"phoneNumber"`
	DateCreated time.Time `json:"dateCreated"`
	IsVerified  bool      `json:"isVerified"`
}

var PaymentInfoService string = "http://localhost:3120"
var connection *amqp.Connection
var emailExchange string = "email"
var emailQueue string = "email"

func init() {
	var err error
	uri := "amqp://guest:guest@localhost:5672/"
	connection, err = amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
}

func NewUser() *User {
	id, _ := uuid.NewV4()
	return &User{
		DateCreated: time.Now(),
		ID:          id.String(),
	}
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password string) error {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PwHash = hpass
	return nil
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
	err = ctx.Database.C("users").Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteUserWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("users").RemoveId(id)
	if err != nil {
		return err
	}
	return nil
}

func FindUsers(ids []string, ctx *DB.Context) []*User {
	adjustedIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		adjustedIDs = append(adjustedIDs, id)
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

func (u *User) SyncWithExistingInvitation(ctx *DB.Context) error {
	test, _ := FindUserByEmail(u.Email, ctx)
	if test != nil {
		if test.IsUserRegistered() {
			return fmt.Errorf("%s", "Email is already registered")
		} else {
			u.ID = test.ID
		}
	}
	return nil
}
