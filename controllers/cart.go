package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/apiresponse"
	"github.com/iskelk/ecommerce-yt/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		productID, ok := productIDFromQuery(c)
		if !ok {
			return
		}

		ctx, cancel := requestContext(c, shortTimeout)
		defer cancel()

		if err := database.AddToCart(ctx, app.ProductCollection, app.UserCollection, productID, userID); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully added to the cart")
	}
}

func (app *Application) RemoveFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		productID, ok := productIDFromQuery(c)
		if !ok {
			return
		}

		ctx, cancel := requestContext(c, shortTimeout)
		defer cancel()

		if err := database.RemoveFromCart(ctx, app.UserCollection, productID, userID); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully removed item from cart")
	}
}

func (app *Application) GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		cart, total, err := database.GetCart(ctx, app.UserCollection, userID)
		if err != nil {
			respondError(c, err)
			return
		}

		respondOK(c, gin.H{
			"total": total,
			"cart":  cart,
		})
	}
}

func (app *Application) CheckoutCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		if err := database.CheckoutCart(ctx, app.UserCollection, userID); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return
		}

		productID, ok := productIDFromQuery(c)
		if !ok {
			return
		}

		ctx, cancel := requestContext(c, shortTimeout)
		defer cancel()

		if err := database.InstantBuy(ctx, app.ProductCollection, app.UserCollection, productID, userID); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully placed the order")
	}
}

func productIDFromQuery(c *gin.Context) (primitive.ObjectID, bool) {
	productQueryID := c.Query("id")
	if productQueryID == "" {
		apiresponse.BadRequest(c, "product id is empty")
		return primitive.NilObjectID, false
	}

	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		apiresponse.BadRequest(c, "product id is not valid")
		return primitive.NilObjectID, false
	}

	return productID, true
}
