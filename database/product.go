package database

import (
	"context"
	"log"
	"regexp"

	"github.com/iskelk/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProduct(ctx context.Context, productCollection *mongo.Collection, product models.Product) error {
	if _, err := productCollection.InsertOne(ctx, product); err != nil {
		log.Println(err)
		return ErrCreateProduct
	}
	return nil
}

func ListProducts(ctx context.Context, productCollection *mongo.Collection) ([]models.Product, error) {
	cursor, err := productCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Println(err)
		return nil, ErrListProducts
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Println(err)
		return nil, ErrListProducts
	}
	return products, nil
}

func SearchProductsByName(ctx context.Context, productCollection *mongo.Collection, query string) ([]models.Product, error) {
	cursor, err := productCollection.Find(ctx, bson.M{
		"product_name": bson.M{"$regex": regexp.QuoteMeta(query), "$options": "i"},
	})
	if err != nil {
		log.Println(err)
		return nil, ErrSearchProducts
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Println(err)
		return nil, ErrSearchProducts
	}
	return products, nil
}
