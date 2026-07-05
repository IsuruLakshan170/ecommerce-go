package controllers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/apiresponse"
	"github.com/iskelk/ecommerce-yt/database"
	"github.com/iskelk/ecommerce-yt/models"
	generate "github.com/iskelk/ecommerce-yt/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *Application) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := requestContext(c, signupTimeout)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			respondBindError(c)
			return
		}

		if validationErr := validateStruct(user); validationErr != nil {
			respondValidationError(c, validationErr)
			return
		}

		exists, err := database.EmailExists(ctx, app.UserCollection, user.Email)
		if err != nil {
			apiresponse.Internal(c, err)
			return
		}
		if exists {
			apiresponse.BadRequest(c, "this email already exists")
			return
		}

		exists, err = database.PhoneExists(ctx, app.UserCollection, user.Phone)
		if err != nil {
			apiresponse.Internal(c, err)
			return
		}
		if exists {
			apiresponse.BadRequest(c, "this phone number already exists")
			return
		}

		hashedPassword, err := hashPassword(*user.Password)
		if err != nil {
			apiresponse.Internal(c, err)
			return
		}
		user.Password = &hashedPassword

		now := time.Now()
		user.Created_At = &now
		user.Updated_At = &now
		user.ID = primitive.NewObjectID()
		userID := user.ID.Hex()
		user.User_ID = &userID

		token, refreshToken, err := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, userID)
		if err != nil {
			apiresponse.Internal(c, err)
			return
		}
		user.Token = &token
		user.Refresh_Token = &refreshToken

		emptyCart := make([]models.ProductUser, 0)
		emptyAddresses := make([]models.Address, 0)
		emptyOrders := make([]models.Order, 0)
		user.User_cart = &emptyCart
		user.User_address_details = &emptyAddresses
		user.User_order_status = &emptyOrders

		if err = database.InsertUser(ctx, app.UserCollection, user); err != nil {
			respondError(c, err)
			return
		}

		respondCreated(c, "successfully signed up")
	}
}

func (app *Application) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			respondBindError(c)
			return
		}

		foundUser, err := database.FindUserByEmail(ctx, app.UserCollection, user.Email)
		if errors.Is(err, database.ErrUserNotFound) {
			apiresponse.Unauthorized(c, apiresponse.MsgLoginIncorrect)
			return
		}
		if err != nil {
			apiresponse.Internal(c, err)
			return
		}

		if !verifyPassword(*user.Password, *foundUser.Password) {
			apiresponse.Unauthorized(c, apiresponse.MsgLoginIncorrect)
			return
		}

		token, refreshToken, err := generate.TokenGenerator(
			*foundUser.Email,
			*foundUser.First_Name,
			*foundUser.Last_Name,
			*foundUser.User_ID,
		)
		if err != nil {
			apiresponse.Internal(c, err)
			return
		}

		if err := database.UpdateUserTokens(ctx, app.UserCollection, *foundUser.User_ID, token, refreshToken); err != nil {
			respondError(c, err)
			return
		}
		foundUser.Token = &token
		foundUser.Refresh_Token = &refreshToken
		foundUser.Password = nil

		respondOK(c, foundUser)
	}
}
