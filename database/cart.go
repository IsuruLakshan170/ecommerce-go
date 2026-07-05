package database

import (
	"context"
	"log"
	"time"

	"github.com/iskelk/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func productToCartItem(productID primitive.ObjectID, product models.Product) models.ProductUser {
	price := float64(0)
	if product.Product_Price != nil {
		price = *product.Product_Price
	}

	productIDStr := productID.Hex()
	item := models.ProductUser{
		Product_ID:       &productIDStr,
		Product_Name:     product.Product_Name,
		Product_Price:    price,
		Product_Category: product.Product_Category,
		Product_Stock:    product.Product_Stock,
		Product_Rating:   product.Product_Rating,
		Product_Reviews:  product.Product_Reviews,
	}
	if product.Product_Image != nil {
		item.Product_Image = *product.Product_Image
	}
	return item
}

func GetCart(ctx context.Context, userCollection *mongo.Collection, userID string) ([]models.ProductUser, float64, error) {
	id, err := parseUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	user, err := FindUserByID(ctx, userCollection, id)
	if err != nil {
		return nil, 0, err
	}

	total, err := cartTotal(ctx, userCollection, id)
	if err != nil {
		return nil, 0, err
	}

	if user.User_cart == nil {
		return []models.ProductUser{}, total, nil
	}
	return *user.User_cart, total, nil
}

func AddToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	var product models.Product
	err := productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	cartItem := productToCartItem(productID, product)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "user_cart", Value: cartItem}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrAddToCart
	}
	return nil
}

func RemoveFromCart(ctx context.Context, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"user_cart": bson.M{"product_id": productID.Hex()}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrRemoveFromCart
	}
	return nil
}

func cartTotal(ctx context.Context, userCollection *mongo.Collection, userID primitive.ObjectID) (float64, error) {
	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user_cart"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	match := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userID}}}}
	group := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$_id"},
		{Key: "total", Value: bson.D{{Key: "$sum", Value: "$user_cart.product_price"}}},
	}}}

	cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{match, unwind, group})
	if err != nil {
		log.Println(err)
		return 0, ErrGetCart
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Println(err)
		return 0, ErrGetCart
	}
	if len(results) == 0 {
		return 0, nil
	}

	switch total := results[0]["total"].(type) {
	case float64:
		return total, nil
	case int32:
		return float64(total), nil
	case int64:
		return float64(total), nil
	default:
		return 0, nil
	}
}

func CheckoutCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	user, err := FindUserByID(ctx, userCollection, id)
	if err != nil {
		return err
	}

	total, err := cartTotal(ctx, userCollection, id)
	if err != nil {
		return err
	}

	orderCart := make([]models.ProductUser, 0)
	if user.User_cart != nil {
		orderCart = append(orderCart, *user.User_cart...)
	}

	orderDetail := models.Order{
		Order_ID:       primitive.NewObjectID(),
		Order_Cart:     orderCart,
		Order_Date:     time.Now(),
		Price:          &total,
		Payment_Method: models.Payment{COD: true},
	}

	emptyCart := make([]models.ProductUser, 0)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{{Key: "user_order_status", Value: orderDetail}}},
		{Key: "$set", Value: bson.D{{Key: "user_cart", Value: emptyCart}}},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCheckoutCart
	}
	return nil
}

func InstantBuy(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := parseUserID(userID)
	if err != nil {
		return err
	}

	var product models.Product
	err = productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	cartItem := productToCartItem(productID, product)
	price := cartItem.Product_Price

	orderDetail := models.Order{
		Order_ID:       primitive.NewObjectID(),
		Order_Cart:     []models.ProductUser{cartItem},
		Order_Date:     time.Now(),
		Price:          &price,
		Payment_Method: models.Payment{COD: true},
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "user_order_status", Value: orderDetail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCheckoutCart
	}
	return nil
}
