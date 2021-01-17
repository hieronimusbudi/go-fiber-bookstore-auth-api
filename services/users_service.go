package services

import (
	resterrors "github.com/hieronimusbudi/go-bookstore-utils/rest_errors"
	"github.com/hieronimusbudi/go-fiber-bookstore-auth-api/domain/users"
	cryptoutils "github.com/hieronimusbudi/go-fiber-bookstore-auth-api/utils/crypto"
	"github.com/hieronimusbudi/go-fiber-bookstore-auth-api/utils/date"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, resterrors.RestErr)
	CreateUser(users.User) (*users.User, resterrors.RestErr)
	UpdateUser(bool, users.User) (*users.User, resterrors.RestErr)
	DeleteUser(int64) resterrors.RestErr
	LoginUser(users.LoginRequest) (*users.User, resterrors.RestErr)
}

func (s *usersService) GetUser(userID int64) (*users.User, resterrors.RestErr) {
	dao := &users.User{ID: userID}
	if err := dao.Get(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, resterrors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date.GetNowDBFormat()
	user.Password = cryptoutils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, resterrors.RestErr) {
	current := &users.User{ID: user.ID}
	if err := current.Get(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}

		if user.LastName != "" {
			current.LastName = user.LastName
		}

		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId int64) resterrors.RestErr {
	dao := &users.User{ID: userId}
	return dao.Delete()
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, resterrors.RestErr) {
	dao := &users.User{
		Email:    request.Email,
		Password: cryptoutils.GetMd5(request.Password),
	}

	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
