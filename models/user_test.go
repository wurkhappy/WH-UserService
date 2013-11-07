package models

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-UserService/DB"
	"testing"
	// "log"
)

func init() {
	config.Test()
	Setup()
	DB.Name = "testdb"
	DB.Setup()
	DB.CreateStatements()
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

		testSaveUser(t)
		testFindUserByEmail(t)
		testFindUserByID(t)
		testDeleteUser(t)
		testFindUsers(t)
		testSyncWithExistingInvitation(t)
		testsendEmail(t)

		DB.DB.Exec("DELETE from wh_user")
	}
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

func testSaveUser(t *testing.T) {
	user := NewUser()
	err := user.Save()

	if err != nil {
		t.Errorf("%s--- error is:%v", "testSaveUser", err)
	}
}

func testFindUserByEmail(t *testing.T) {
	user := NewUser()
	user.Email = "test@test.com"
	err := user.Save()
	if err != nil {
		t.Errorf("%s--- error saving", "testFindUserByEmail")
	}

	u, err := FindUserByEmail(user.Email)
	if err != nil {
		t.Errorf("testFindUserByEmail--- error finding user %v", err)
	}

	if u == nil {
		t.Errorf("%s--- user was not found", "testFindUserByEmail")
	}

	_, err = FindUserByEmail("badrequest@bad.com")
	if err == nil {
		t.Errorf("%s--- DB returned a bad request", "testFindUserByEmail")
	}
}

func testFindUserByID(t *testing.T) {
	user := NewUser()
	user.Save()

	u, err := FindUserByID(user.ID)
	if err != nil {
		t.Errorf("testFindUserByID--- error finding user %v", err)
	}

	if u == nil {
		t.Errorf("%s--- user was not found", "testFindUserByID")
	}

	_, err = FindUserByID("invalidID")
	if err == nil {
		t.Errorf("%s--- DB returned a bad request", "testFindUserByID")
	}
}

func testDeleteUser(t *testing.T) {
	user := NewUser()
	user.Save()
	err := DeleteUserWithID(user.ID)

	if err != nil {
		t.Errorf("%s--- error is:%v", "testDeleteUser", err)
	}

	u, err := FindUserByID(user.ID)
	if u != nil {
		t.Errorf("%s--- user was found", "testDeleteUser")
	}

	err = DeleteUserWithID("invalid-id")
	if err == nil {
		t.Errorf("%s--- DB deleted with invalid id", "testDeleteUser")
	}
}

func testFindUsers(t *testing.T) {
	user1 := NewUser()
	user2 := NewUser()

	user1.Save()
	user2.Save()

	users := FindUsers([]string{user1.ID, user2.ID})
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

func testSyncWithExistingInvitation(t *testing.T) {
	user := NewUser()
	user.Email = "testSyncWithExistingInvitation@wh.com"
	user.Save()

	user2 := new(User)
	user2.Email = user.Email
	err := user2.SyncWithExistingInvitation()
	if err != nil {
		t.Errorf("%s--- didn't sync user", "testSyncWithExistingInvitation")
	}
	if user2.ID != user.ID {
		t.Errorf("%s--- didn't sync ID", "testSyncWithExistingInvitation")
	}

	user.SetPassword("password")
	user.Save()

	user3 := new(User)
	user3.Email = user.Email
	err = user3.SyncWithExistingInvitation()
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
