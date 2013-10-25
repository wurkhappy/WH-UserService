package models

import (
	"encoding/json"
	"github.com/wurkhappy/WH-UserService/DB"
	"labix.org/v2/mgo"
	"testing"
)

func init() {
	emailExchange = "test"
	emailQueue = "test"
}

func TestUnitTests(t *testing.T) {
	testNewUser(t)
	testSetPassword(t)
	testPasswordIsValid(t)
	testValidateNewPassword(t)
	testIsUserRegistered(t)
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
		testFindUsers(t, ctx)
		testSyncWithExistingInvitation(t, ctx)
		testsendEmail(t)
	}
}

//helpers
func ClearDB(ctx *DB.Context) {
	ctx.Database.C("users").DropCollection()
}

func testInit(t *testing.T) {
	if connection == nil {
		t.Errorf("%s--- rabbitmq connection wasn't initialized", "testInit")
	}
}

//tests
func testNewUser(t *testing.T) {
	user := NewUser()

	if user.ID == "" {
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

	err := user.SetPassword("")
	if err == nil {
		t.Errorf("%s--- does not return error on invalid input", "testSetPassword")
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

	_, err = FindUserByEmail("badrequest@bad.com", ctx)
	if err == nil {
		t.Errorf("%s--- DB returned a bad request", "testFindUserByEmail")
	}
}

func testFindUserByID(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	user.SaveUserWithCtx(ctx)

	u, _ := FindUserByID(user.ID, ctx)

	if len(u.ID) == 0 {
		t.Errorf("%s--- user was not found", "testFindUserByID")
	}
}

func testDeleteUser(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	user.SaveUserWithCtx(ctx)
	err := DeleteUserWithID(user.ID, ctx)

	if err != nil {
		t.Errorf("%s--- error is:%v", "testDeleteUser", err)
	}

	u, err := FindUserByID(user.ID, ctx)
	if u != nil {
		t.Errorf("%s--- user was found", "testDeleteUser")
	}

	err = DeleteUserWithID("invalid-id", ctx)
	if err == nil {
		t.Errorf("%s--- DB deleted with invalid id", "testDeleteUser")
	}
}

func testFindUsers(t *testing.T, ctx *DB.Context) {
	user1 := NewUser()
	user2 := NewUser()

	user1.SaveUserWithCtx(ctx)
	user2.SaveUserWithCtx(ctx)

	users := FindUsers([]string{user1.ID, user2.ID}, ctx)
	if len(users) != 2 {
		t.Errorf("%s--- all users were not found", "testFindUsers")
	}
}

func testPasswordIsValid(t *testing.T) {
	user := new(User)
	password := "password"
	user.SetPassword(password)

	if !user.PasswordIsValid(password) {
		t.Errorf("%s--- invalid password", "testPasswordIsValid")
	}

	if user.PasswordIsValid("invalid") {
		t.Errorf("%s--- invalid password passed", "testPasswordIsValid")
	}
}

func testValidateNewPassword(t *testing.T) {
	user := new(User)
	err := user.ValidateNewPassword("123456")
	if err != nil {
		t.Errorf("%s--- valid password failing", "testValidateNewPassword")
	}

	err = user.ValidateNewPassword("12345")
	if err == nil {
		t.Errorf("%s--- invalid password validating", "testValidateNewPassword")
	}
}

func testIsUserRegistered(t *testing.T) {
	user := new(User)
	if user.IsUserRegistered() {
		t.Errorf("%s--- unregistered user is returning as registered", "testIsUserRegistered")
	}
	user.PwHash = []byte("some-string")
	if !user.IsUserRegistered() {
		t.Errorf("%s--- registered user is returning as unregistered", "testIsUserRegistered")
	}
}

func testSyncWithExistingInvitation(t *testing.T, ctx *DB.Context) {
	user := NewUser()
	user.Email = "testSyncWithExistingInvitation@wh.com"
	user.SaveUserWithCtx(ctx)

	user2 := new(User)
	user2.Email = user.Email
	err := user2.SyncWithExistingInvitation(ctx)
	if err != nil {
		t.Errorf("%s--- didn't sync user", "testSyncWithExistingInvitation")
	}
	if user2.ID != user.ID {
		t.Errorf("%s--- didn't sync ID", "testSyncWithExistingInvitation")
	}

	user.SetPassword("password")
	user.SaveUserWithCtx(ctx)

	user3 := new(User)
	user3.Email = user.Email
	err = user3.SyncWithExistingInvitation(ctx)
	if err == nil {
		t.Errorf("%s--- failed to reject registered user", "testSyncWithExistingInvitation")
	}
}

func testsendEmail(t *testing.T) {
	path := "test"
	u := NewUser()
	message := map[string]interface{}{
		"Body": map[string]interface{}{
			"user": u,
		},
	}
	body, _ := json.Marshal(message)
	err := sendEmail(path, body)
	if err != nil {
		t.Errorf("%s--- error sending email", "testsendEmail")
	}

}
