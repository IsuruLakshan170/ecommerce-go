package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/database"
	"github.com/iskelk/ecommerce-yt/models"
)

func (app *Application) AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		var address models.Address
		if err := c.BindJSON(&address); err != nil {
			respondBindError(c)
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		if err := database.AddAddress(ctx, app.UserCollection, userID, address); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully added address")
	}
}

func (app *Application) EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		var address models.Address
		if err := c.BindJSON(&address); err != nil {
			respondBindError(c)
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		if err := database.UpdateHomeAddress(ctx, app.UserCollection, userID, address); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully updated home address")
	}
}

func (app *Application) EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		var address models.Address
		if err := c.BindJSON(&address); err != nil {
			respondBindError(c)
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		if err := database.UpdateWorkAddress(ctx, app.UserCollection, userID, address); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully updated work address")
	}
}

func (app *Application) DeleteAddresses() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		if err := database.DeleteAddresses(ctx, app.UserCollection, userID); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully deleted addresses")
	}
}
