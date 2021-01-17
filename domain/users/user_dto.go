package users

import (
	"strings"

	resterrors "github.com/hieronimusbudi/go-bookstore-utils/rest_errors"
)

const (
	StatusActive = "active"
)

type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Validate func
func (user *User) Validate() resterrors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return resterrors.NewBadRequestError("invalid email/password")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return resterrors.NewBadRequestError("invalid email/password")
	}

	return nil
}
