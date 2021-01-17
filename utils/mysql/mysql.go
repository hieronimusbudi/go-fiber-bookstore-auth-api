package mysqlutils

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	resterrors "github.com/hieronimusbudi/go-bookstore-utils/rest_errors"
)

const (
	ErrorNoRows = "no  rows in result set"
)

func ParseError(err error) resterrors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return resterrors.NewNotFoundError("no record matching given id")
		}
		return resterrors.NewInternalServerError("error parsing database response", err)
	}

	switch sqlErr.Number {
	case 1062:
		return resterrors.NewBadRequestError("invalid data")
	}

	return resterrors.NewInternalServerError("database error", err)
}
