package schema

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type UserNotifications struct {
	ID             int    `json:"id"`
	UserID         string `json:"user_id"`
	NotificationID string `json:"notification_id"`
	Email          string `json:"email"`
}

func (u *User) HashPassword() error {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) Authorized(user User) bool {
	hashFromDatabase := user.Password
	err := bcrypt.CompareHashAndPassword([]byte(hashFromDatabase), []byte(u.Password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (u User) Validate() error {
	var errStr string

	reqStrFields := make(map[string]string)
	reqStrFields["username"] = u.UserName
	reqStrFields["password"] = u.Password
	reqStrFields["email"] = u.Email

	for k, v := range reqStrFields {
		nes := nonEmptyString(k, v)
		if nes != "" {
			errStr += nes
		}
	}

	if errStr != "" {
		return errors.New(strings.TrimSuffix(errStr, " "))
	}
	return nil
}

func nonEmptyString(key, value string) string {
	if value == "" {
		return requiredFieldMessage(key)
	}
	return ""
}

func requiredFieldMessage(field string) string {
	return fmt.Sprintf("%s is a required field. ", field)
}
