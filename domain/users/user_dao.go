package users

import (
	"github.com/hieronimusbudi/go-bookstore-utils/logger"
	resterrors "github.com/hieronimusbudi/go-bookstore-utils/rest_errors"
	"github.com/hieronimusbudi/go-fiber-bookstore-auth-api/datasouces/mysql/users_db"
	mysqlutils "github.com/hieronimusbudi/go-fiber-bookstore-auth-api/utils/mysql"

	"errors"
	"strings"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, status, password) VALUES(?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?"
)

func (user *User) Get() resterrors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		// logger.Error("error when trying to prepare get user statement", err)
		return resterrors.NewInternalServerError("error when tying to get user", err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(&user.ID)

	if getErr := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		// logger.Error("error when trying to get user by id", getErr)
		return resterrors.NewInternalServerError("error when tying to get user", getErr)
	}
	return nil
}

func (user *User) Save() resterrors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		// logger.Error("error when trying to prepare save user statement", err)
		return resterrors.NewInternalServerError("error when tying to save user 1", err)
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password)
	if saveErr != nil {
		// logger.Error("error when trying to save user", saveErr)
		return resterrors.NewInternalServerError("error when tying to save user", saveErr)
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		// logger.Error("error when trying to get last insert id after creating a new user", err)
		return resterrors.NewInternalServerError("error when tying to save user", err)
	}
	user.ID = userID

	return nil
}

func (user *User) Update() resterrors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		// logger.Error("error when trying to prepare update user statement", err)
		return resterrors.NewInternalServerError("error when tying to update user", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		// logger.Error("error when trying to update user", err)
		return resterrors.NewInternalServerError("error when tying to update user", err)
	}
	return nil
}

func (user *User) Delete() resterrors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		// logger.Error("error when trying to prepare delete user statement", err)
		return resterrors.NewInternalServerError("error when tying to update user", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.ID); err != nil {
		// logger.Error("error when trying to delete user", err)
		return resterrors.NewInternalServerError("error when tying to save user", err)
	}
	return nil
}

func (user *User) FindByEmailAndPassword() resterrors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return resterrors.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysqlutils.ErrorNoRows) {
			return resterrors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return resterrors.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	return nil
}
