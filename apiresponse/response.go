package apiresponse

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MsgInternalError  = "something went wrong, please try again later"
	MsgInvalidBody    = "invalid request body"
	MsgInvalidToken   = "invalid token"
	MsgTokenExpired   = "token is expired"
	MsgNoToken        = "no authorization token provided"
	MsgLoginIncorrect = "login or password is incorrect"
)

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func Internal(c *gin.Context, err error) {
	if err != nil {
		log.Println(err)
	}
	Error(c, http.StatusInternalServerError, MsgInternalError)
}

func Message(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}
