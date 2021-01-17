package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	jwtutils "github.com/hieronimusbudi/go-bookstore-utils/jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	resterrors "github.com/hieronimusbudi/go-bookstore-utils/rest_errors"
	"github.com/hieronimusbudi/go-fiber-bookstore-auth-api/domain/users"
	"github.com/hieronimusbudi/go-fiber-bookstore-auth-api/services"
	"github.com/valyala/fasthttp"
)

var (
	// jwtSecret     = os.Getenv("JWT_SECRET")
	// jwtCookieName = os.Getenv("JWT_COOKIE_NAME")
	jwtSecret     = "secret"
	jwtCookieName = "token::jwt"
)

func getUserId(userIDParam string) (int64, resterrors.RestErr) {
	userID, userErr := strconv.ParseInt(userIDParam, 10, 64)
	if userErr != nil {
		return 0, resterrors.NewBadRequestError("user id should be a number")
	}
	return userID, nil
}

func Create(c *fiber.Ctx) error {
	user := new(users.User)
	if err := c.BodyParser(user); err != nil {
		restErr := resterrors.NewBadRequestError("invalid json body")
		return c.Status(restErr.Status()).JSON(restErr)
	}

	result, saveErr := services.UsersService.CreateUser(*user)
	if saveErr != nil {
		return c.Status(saveErr.Status()).JSON(saveErr)
	}

	token, tokenErr := jwtutils.GenerateToken(&jwtutils.UserPayload{
		Id:        result.ID,
		Email:     result.Email,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Status:    result.Status,
	}, jwtSecret)
	if tokenErr != nil {
		return c.Status(tokenErr.Status()).JSON(tokenErr)
	}

	cookie := new(fasthttp.Cookie)
	cookie.SetKey(jwtCookieName)
	cookie.SetValue(token)
	cookie.SetExpire(time.Now().Add(time.Hour * time.Duration(1)))
	c.Response().Header.SetCookie(cookie)

	return c.Status(http.StatusCreated).JSON(result)
}

func Get(c *fiber.Ctx) error {
	userId, idErr := getUserId(c.Params("user_id"))
	if idErr != nil {
		return c.Status(idErr.Status()).JSON(idErr)

	}
	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		return c.Status(getErr.Status()).JSON(getErr)
	}

	tokenClaims, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)
	if !ok {
		restJwtErr := resterrors.NewUnauthorizedError("Token claims not exists")
		return c.Status(restJwtErr.Status()).JSON(restJwtErr)
	}

	log.Printf("check %v\n", tokenClaims)
	if int64(tokenClaims["id"].(float64)) == user.ID {
		return c.Status(http.StatusOK).JSON(user.Marshall(false))
	}
	return c.Status(http.StatusOK).JSON(user.Marshall(true))
}

func Login(c *fiber.Ctx) error {
	request := new(users.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		restErr := resterrors.NewBadRequestError("invalid json body")
		return c.Status(restErr.Status()).JSON(restErr)
	}

	user, err := services.UsersService.LoginUser(*request)
	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	// create jwt token
	token, tokenErr := jwtutils.GenerateToken(&jwtutils.UserPayload{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    user.Status,
	}, jwtSecret)
	if tokenErr != nil {
		return c.Status(tokenErr.Status()).JSON(tokenErr)
	}

	cookie := new(fasthttp.Cookie)
	cookie.SetKey(jwtCookieName)
	cookie.SetValue(token)
	cookie.SetExpire(time.Now().Add(time.Hour * time.Duration(1)))
	c.Response().Header.SetCookie(cookie)

	return c.Status(http.StatusCreated).JSON(user.Marshall(true))
}

// func Update(c *gin.Context) {
// 	// Get token from cookie
// 	jwtCookie, jwtErr := c.Request.Cookie(jwtCookieName)
// 	if jwtErr != nil {
// 		restJwtErr := resterrors.NewUnauthorizedError(jwtErr.Error())
// 		c.JSON(restJwtErr.Status(), restJwtErr)
// 		return
// 	}

// 	// Validate token
// 	_, tokenErr := myjwt.ValidateToken(jwtCookie.Value, jwtSecret)
// 	if tokenErr != nil {
// 		c.JSON(tokenErr.Status(), tokenErr)
// 		return
// 	}

// 	userId, idErr := getUserId(c.Param("user_id"))
// 	if idErr != nil {
// 		c.JSON(idErr.Status(), idErr)
// 		return
// 	}

// 	var user users.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		restErr := resterrors.NewBadRequestError("invalid json body")
// 		c.JSON(restErr.Status(), restErr)
// 		return
// 	}

// 	user.ID = userId

// 	isPartial := c.Request.Method == http.MethodPatch

// 	result, err := services.UsersService.UpdateUser(isPartial, user)
// 	if err != nil {
// 		c.JSON(err.Status(), err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
// }

// func Delete(c *gin.Context) {
// 	// Get token from cookie
// 	jwtCookie, jwtErr := c.Request.Cookie(jwtCookieName)
// 	if jwtErr != nil {
// 		restJwtErr := resterrors.NewUnauthorizedError(jwtErr.Error())
// 		c.JSON(restJwtErr.Status(), restJwtErr)
// 		return
// 	}

// 	// Validate token
// 	_, tokenErr := myjwt.ValidateToken(jwtCookie.Value, jwtSecret)
// 	if tokenErr != nil {
// 		c.JSON(tokenErr.Status(), tokenErr)
// 		return
// 	}

// 	userID, idErr := getUserId(c.Param("user_id"))
// 	if idErr != nil {
// 		c.JSON(idErr.Status(), idErr)
// 		return
// 	}

// 	if err := services.UsersService.DeleteUser(userID); err != nil {
// 		c.JSON(err.Status(), err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
// }
