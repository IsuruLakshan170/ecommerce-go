package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/controllers"
)

func PublicRoutes(router *gin.Engine, app *controllers.Application) {
	router.POST("/users/signup", app.SignUp())
	router.POST("/users/login", app.Login())
	router.GET("/users/productview", app.ListProducts())
	router.GET("/users/search", app.SearchProducts())
}

func ProtectedRoutes(router *gin.RouterGroup, app *controllers.Application) {
	router.POST("/admin/addproduct", app.AddProduct())

	router.POST("/addtocart", app.AddToCart())
	router.DELETE("/removeitem", app.RemoveFromCart())
	router.GET("/usercart", app.GetCart())
	router.POST("/cartcheckout", app.CheckoutCart())
	router.POST("/instantbuy", app.InstantBuy())

	router.POST("/users/addaddress", app.AddAddress())
	router.PUT("/users/edithomeaddress", app.EditHomeAddress())
	router.PUT("/users/editworkaddress", app.EditWorkAddress())
	router.DELETE("/users/deleteaddress", app.DeleteAddresses())
}
