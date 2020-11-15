package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"goacs/http/request"
	"goacs/http/response"
	"goacs/lib"
	"goacs/models/user"
	"goacs/repository"
	"goacs/repository/mysql"
	"log"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  user.User `json:"user"`
	Token string    `json:"token"`
}

func Login(ctx *gin.Context) {
	var loginRequest LoginRequest

	err := json.NewDecoder(ctx.Request.Body).Decode(&loginRequest)

	if err != nil {
		log.Println("Error in req ", err)
	}

	validator := request.NewApiValidator(ctx, loginRequest)
	verr := validator.Validate()

	if verr != nil {
		response.ResponseValidationErrors(ctx, validator)
		return
	}

	userRepository := mysql.NewUserRepository(repository.GetConnection())
	userByAuthData, err := userRepository.GetUserByAuthData(loginRequest.Username, loginRequest.Password)

	if err != nil {
		log.Println("Cannot find userByAuthData", err.Error())
		response.ResponseError(ctx, 404, "Cannot find userByAuthData", err)
		return
	}

	loginResponse := LoginResponse{
		User:  userByAuthData,
		Token: NewTokenForUser(userByAuthData),
	}

	response.ResponseData(ctx, loginResponse)
}

func NewTokenForUser(user user.User) string {
	env := new(lib.Env)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 120).Unix(),
		Subject:   user.Uuid,
		Issuer:    "user",
	})

	tokenString, err := token.SignedString([]byte(env.Get("JWT_SECRET", "")))
	if err != nil {
		log.Println("Error while generating token ", err)
	}
	return tokenString
}
