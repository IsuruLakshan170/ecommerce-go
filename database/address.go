package database

import (
	"context"
	"log"
	"strconv"

	"github.com/iskelk/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress(ctx context.Context, userCollection *mongo.Collection, userID string, address models.Address) error {
	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	user, err := FindUserByID(ctx, userCollection, id)
	if err != nil {
		return err
	}

	addressCount := 0
	if user.User_address_details != nil {
		addressCount = len(*user.User_address_details)
	}
	if addressCount >= 2 {
		return ErrTooManyAddresses
	}

	address.Address_ID = primitive.NewObjectID()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "user_address_details", Value: address}}}}
	if _, err = userCollection.UpdateOne(ctx, filter, update); err != nil {
		log.Println(err)
		return ErrAddAddress
	}
	return nil
}

func UpdateHomeAddress(ctx context.Context, userCollection *mongo.Collection, userID string, address models.Address) error {
	return updateAddressAtIndex(ctx, userCollection, userID, 0, address)
}

func UpdateWorkAddress(ctx context.Context, userCollection *mongo.Collection, userID string, address models.Address) error {
	return updateAddressAtIndex(ctx, userCollection, userID, 1, address)
}

func updateAddressAtIndex(ctx context.Context, userCollection *mongo.Collection, userID string, index int, address models.Address) error {
	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	prefix := "user_address_details." + strconv.Itoa(index) + "."
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: prefix + "address_street", Value: address.Address_Street},
		{Key: prefix + "address_city", Value: address.Address_City},
		{Key: prefix + "address_state", Value: address.Address_State},
		{Key: prefix + "address_zip", Value: address.Address_Zip},
		{Key: prefix + "address_country", Value: address.Address_Country},
	}}}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrUpdateAddress
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

func DeleteAddresses(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	emptyAddresses := make([]models.Address, 0)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "user_address_details", Value: emptyAddresses}}}}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrDeleteAddress
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
