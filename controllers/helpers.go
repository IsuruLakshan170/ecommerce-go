package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/apiresponse"
	"github.com/iskelk/ecommerce-yt/database"
)

const (
	defaultTimeout = 100 * time.Second
	shortTimeout   = 5 * time.Second
	signupTimeout  = 10 * time.Second
)

func requestContext(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), timeout)
}

func respondBindError(c *gin.Context) {
	apiresponse.BadRequest(c, apiresponse.MsgInvalidBody)
}

func respondValidationError(c *gin.Context, err error) {
	apiresponse.BadRequest(c, err.Error())
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, database.ErrUserIdIsNotValid):
		apiresponse.BadRequest(c, "user id is not valid")
	case errors.Is(err, database.ErrUserNotFound):
		apiresponse.NotFound(c, "user not found")
	case errors.Is(err, database.ErrCantFindProduct):
		apiresponse.NotFound(c, "product not found")
	case errors.Is(err, database.ErrTooManyAddresses):
		apiresponse.BadRequest(c, err.Error())
	case errors.Is(err, database.ErrAddToCart),
		errors.Is(err, database.ErrRemoveFromCart),
		errors.Is(err, database.ErrGetCart),
		errors.Is(err, database.ErrCheckoutCart),
		errors.Is(err, database.ErrAddAddress),
		errors.Is(err, database.ErrUpdateAddress),
		errors.Is(err, database.ErrDeleteAddress),
		errors.Is(err, database.ErrCreateUser),
		errors.Is(err, database.ErrCreateProduct),
		errors.Is(err, database.ErrListProducts),
		errors.Is(err, database.ErrSearchProducts),
		errors.Is(err, database.ErrUpdateTokens),
		errors.Is(err, database.ErrUpdateUser):
		apiresponse.Error(c, http.StatusInternalServerError, err.Error())
	default:
		apiresponse.Internal(c, err)
	}
}

func respondSuccess(c *gin.Context, status int, message string) {
	apiresponse.Message(c, status, message)
}

func respondOK(c *gin.Context, payload any) {
	c.JSON(http.StatusOK, payload)
}

func respondCreated(c *gin.Context, message string) {
	apiresponse.Message(c, http.StatusCreated, message)
}
