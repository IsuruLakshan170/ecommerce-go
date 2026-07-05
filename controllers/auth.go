package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/apiresponse"
	"golang.org/x/crypto/bcrypt"
)

func authenticatedUserID(c *gin.Context) (string, bool) {
	uid, exists := c.Get("uid")
	if !exists {
		apiresponse.Unauthorized(c, "user not authenticated")
		return "", false
	}

	userID, ok := uid.(string)
	if !ok || userID == "" {
		apiresponse.Unauthorized(c, "invalid user session")
		return "", false
	}

	return userID, true
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func verifyPassword(plainPassword, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}
