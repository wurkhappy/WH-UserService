package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-UserService/DB"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
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
	IsVerified  bool          `json:"isVerified"`
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
		panic(err)
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

func (u *User) AddToPaymentProcessor() {
	client := &http.Client{}
	r, _ := http.NewRequest("POST", "http://localhost:3120/user/"+u.ID.Hex(), nil)
	_, err := client.Do(r)
	if err != nil {
	}
}

func (u *User) SendVerificationEmail() {
	message := map[string]interface{}{
		"Body": u,
	}
	uri := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/user/verify")
	publisher.Publish(body, true)
}

func (u *User) SendForgotPasswordEmail() {
	message := map[string]interface{}{
		"Body": map[string]interface{}{
			"user":     u,
		},
	}
	uri := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/user/password/forgot")
	publisher.Publish(body, true)
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
