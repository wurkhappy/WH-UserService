package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/dchest/uniuri"
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
	Fingerprint string        `json:"-"`
	FirstName   string        `json:"firstName"`
	LastName    string        `json:"lastName"`
	Email       string        `json:"email"`
	PwHash      []byte        `json:"-"`
	AvatarURL   string        `json:"avatarURL"`
	PhoneNumber string        `json:"phoneNumber"`
	DateCreated time.Time     `json:"dateCreated"`
	IsVerified  bool          `json:"isVerified"`
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
		Fingerprint: uniuri.New(),
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

func (u *User) AddToPaymentProcessor() {
	client := &http.Client{}
	r, _ := http.NewRequest("POST", "http://localhost:3120/user/"+u.ID.Hex(), nil)
	_, err := client.Do(r)
	if err != nil {
	}
}

func (u *User) CreateSignature(path string) string {
	mac := hmac.New(sha256.New, []byte(u.Fingerprint))
	mac.Write([]byte(path))
	log.Print(path)
	return hex.EncodeToString(mac.Sum(nil))
}

func (u *User) VerifySignature(path, signature string) bool {
	mac := hmac.New(sha256.New, []byte(u.Fingerprint))
	mac.Write([]byte(path))
	log.Print(path)
	log.Print(signature)

	sigMAC, _ := hex.DecodeString(signature)
	if !hmac.Equal(sigMAC, mac.Sum(nil)) {
		return false
	}
	return true
}

func (u *User) SendVerificationEmail() {
	message := map[string]interface{}{
		"Body":   u,
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

func (u *User) SendNewPasswordEmail() {
	message := map[string]interface{}{
		"Body":   u,
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
