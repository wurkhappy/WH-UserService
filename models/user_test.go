package models

import (
	"github.com/wurkhappy/WH-UserService/DB"
	"labix.org/v2/mgo"
	"testing"
	"time"
)

func TestUnitTests(t *testing.T) {
	testNewUser(t)
	testSetPassword(t)
	testPasswordIsValid(t)
	testVerifySignature(t)
}

func TestIntegrationTests(t *testing.T) {
	if !testing.Short() {
		DB.Session, _ = mgo.Dial(DB.Config["DBURL"])
		ctx := &DB.Context{
			Database: DB.Session.Clone().DB("TestUserDB"),
		}
		defer ctx.Close()
		defer ClearDB(ctx)

		testSaveUser(t, ctx)
		testFindUserByEmail(t, ctx)
		testFindUserByID(t, ctx)
		testDeleteUser(t, ctx)

	}
}

//helpers
func ClearDB(ctx *DB.Context) {
	ctx.Database.C("users").DropCollection()
}

//tests
func testNewUser(t *testing.T) {
	user := NewUser()

	if !user.ID.Valid() {
		t.Errorf("%s--- user id not valid", "testNewUser")
	}
}

func testSetPassword(t *testing.T) {
	user := new(User)
	password := "password"
	user.SetPassword(password)

	if output, input := user.PwHash, password; string(output) == input {
		t.Errorf("%s--- input:%s, output: %s", "testSetPassword", input, output)
	}
}

func testSaveUser(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	err := user.SaveUserWithCtx(ctx)

	if err != nil {
		t.Errorf("%s--- error is:%v", "testSaveUser", err)
	}
}

func testFindUserByEmail(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	user.Email = "test@test.com"
	err := user.SaveUserWithCtx(ctx)
	if err != nil {
		t.Errorf("%s--- error saving", "testFindUserByEmail")
	}

	u, err := FindUserByEmail(user.Email, ctx)
	if err != nil {
		t.Errorf("testFindUserByEmail--- error finding user %v", err)
	}

	if u == nil {
		t.Errorf("%s--- user was not found", "testFindUserByEmail")
	}
}

func testFindUserByID(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	user.SaveUserWithCtx(ctx)

	u, _ := FindUserByID(user.ID.Hex(), ctx)

	if len(u.ID.Hex()) == 0 {
		t.Errorf("%s--- user was not found", "testFindUserByID")
	}
}

func testDeleteUser(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	user.SaveUserWithCtx(ctx)
	err := DeleteUserWithID(user.ID.Hex(), ctx)

	if err != nil {
		t.Errorf("%s--- error is:%v", "testDeleteUser", err)
	}

	u, err := FindUserByID(user.ID.Hex(), ctx)
	if u != nil {
		t.Errorf("%s--- user was found", "testDeleteUser")
	}
}

func testPasswordIsValid(t *testing.T) {
	user := new(User)
	password := "password"
	user.SetPassword(password)

	user.PasswordIsValid(password)

	if !user.PasswordIsValid(password) {
		t.Errorf("%s--- invalid password", "testPasswordIsValid")
	}
}

func testVerifySignature(t *testing.T) {
	path := "/test/path"
	expiration := int(time.Now().Unix())
	user := NewUser()
	signature := user.CreateSignature(path, expiration)
	if !user.VerifySignature(path, expiration, signature) {
		t.Error("Signature not being verified")
	}
}
