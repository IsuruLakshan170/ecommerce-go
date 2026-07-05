package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	UserCollection    *mongo.Collection
	ProductCollection *mongo.Collection
}

func NewApplication(userCollection, productCollection *mongo.Collection) *Application {
	return &Application{
		UserCollection:    userCollection,
		ProductCollection: productCollection,
	}
}
