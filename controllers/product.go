package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/apiresponse"
	"github.com/iskelk/ecommerce-yt/database"
	"github.com/iskelk/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *Application) AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			respondBindError(c)
			return
		}

		product.ID = primitive.NewObjectID()
		productID := product.ID.Hex()
		product.Product_ID = &productID

		if err := database.InsertProduct(ctx, app.ProductCollection, product); err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, 200, "successfully added product")
	}
}

func (app *Application) ListProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		products, err := database.ListProducts(ctx, app.ProductCollection)
		if err != nil {
			respondError(c, err)
			return
		}

		respondOK(c, products)
	}
}

func (app *Application) SearchProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParam := c.Query("name")
		if queryParam == "" {
			apiresponse.BadRequest(c, "invalid search query")
			return
		}

		ctx, cancel := requestContext(c, defaultTimeout)
		defer cancel()

		products, err := database.SearchProductsByName(ctx, app.ProductCollection, queryParam)
		if err != nil {
			respondError(c, err)
			return
		}

		respondOK(c, products)
	}
}
