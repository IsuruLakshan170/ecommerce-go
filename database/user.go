package database

import (
	"context"
	"fmt"
	"log"

	"github.com/iskelk/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func parseUserID(userID string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, ErrUserIdIsNotValid
	}
	return id, nil
}

func EmailExists(ctx context.Context, userCollection *mongo.Collection, email *string) (bool, error) {
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, fmt.Errorf("count users by email: %w", err)
	}
	return count > 0, nil
}

func PhoneExists(ctx context.Context, userCollection *mongo.Collection, phone *string) (bool, error) {
	count, err := userCollection.CountDocuments(ctx, bson.M{"phone": phone})
	if err != nil {
		return false, fmt.Errorf("count users by phone: %w", err)
	}
	return count > 0, nil
}

func InsertUser(ctx context.Context, userCollection *mongo.Collection, user models.User) error {
	if _, err := userCollection.InsertOne(ctx, user); err != nil {
		log.Println(err)
		return ErrCreateUser
	}
	return nil
}

func FindUserByEmail(ctx context.Context, userCollection *mongo.Collection, email *string) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, ErrUserNotFound
		}
		log.Println(err)
		return user, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

func FindUserByID(ctx context.Context, userCollection *mongo.Collection, userID primitive.ObjectID) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: userID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, ErrUserNotFound
		}
		log.Println(err)
		return user, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

func SetUserAdmin(ctx context.Context, userCollection *mongo.Collection, userID string, isAdmin bool) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"is_admin": isAdmin}}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrUpdateUser
	}
	return nil
}

func UpdateUserTokens(ctx context.Context, userCollection *mongo.Collection, userID, token, refreshToken string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
		},
	}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrUpdateTokens
	}
	return nil
}
