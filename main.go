package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/controllers"
	"github.com/iskelk/ecommerce-yt/database"
	"github.com/iskelk/ecommerce-yt/middleware"
	"github.com/iskelk/ecommerce-yt/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		if err := client.Disconnect(disconnectCtx); err != nil {
			log.Println("mongo disconnect:", err)
		}
	}()

	fmt.Println("Connected to MongoDB")

	app := controllers.NewApplication(
		database.UsersCollection(client),
		database.ProductsCollection(client),
	)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	routes.PublicRoutes(router, app)

	authorized := router.Group("/")
	authorized.Use(middleware.Authentication())
	routes.ProtectedRoutes(authorized, app)

	log.Fatal(router.Run(":" + port))
}
