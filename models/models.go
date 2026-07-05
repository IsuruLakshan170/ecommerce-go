package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                   primitive.ObjectID `bson:"_id"`
	First_Name           *string            `json:"first_name" bson:"first_name" validate:"required,min=2,max=30"`
	Last_Name            *string            `json:"last_name" bson:"last_name" validate:"required,min=2,max=30"`
	Password             *string            `json:"password" bson:"password" validate:"required,min=6"`
	Email                *string            `json:"email" bson:"email" validate:"email,required"`
	Phone                *string            `json:"phone" bson:"phone" validate:"required"`
	Token                *string            `json:"token" bson:"token"`
	Refresh_Token        *string            `json:"refresh_token" bson:"refresh_token"`
	Created_At           *time.Time         `json:"created_at" bson:"created_at"`
	Updated_At           *time.Time         `json:"updated_at" bson:"updated_at"`
	User_ID              *string            `json:"user_id" bson:"user_id"`
	User_cart            *[]ProductUser     `json:"user_cart" bson:"user_cart"`
	User_address_details *[]Address         `json:"user_address_details" bson:"user_address_details"`
	User_order_status    *[]Order           `json:"user_order_status" bson:"user_order_status"`
}

type Product struct {
	ID                  primitive.ObjectID `bson:"_id"`
	Product_ID          *string            `json:"product_id" bson:"product_id"`
	Product_Name        string             `json:"product_name" bson:"product_name"`
	Product_Description *string            `json:"product_description" bson:"product_description"`
	Product_Price       *float64           `json:"product_price" bson:"product_price"`
	Product_Image       *string            `json:"product_image" bson:"product_image"`
	Product_Category    string             `json:"product_category" bson:"product_category"`
	Product_Stock       int                `json:"product_stock" bson:"product_stock"`
	Product_Rating      float64            `json:"product_rating" bson:"product_rating"`
	Product_Reviews     []string           `json:"product_reviews" bson:"product_reviews"`
}

type ProductUser struct {
	Product_ID       *string  `json:"product_id" bson:"product_id"`
	Product_Name     string   `json:"product_name" bson:"product_name"`
	Product_Price    float64  `json:"product_price" bson:"product_price"`
	Product_Image    string   `json:"product_image" bson:"product_image"`
	Product_Category string   `json:"product_category" bson:"product_category"`
	Product_Stock    int      `json:"product_stock" bson:"product_stock"`
	Product_Rating   float64  `json:"product_rating" bson:"product_rating"`
	Product_Reviews  []string `json:"product_reviews" bson:"product_reviews"`
}

type Order struct {
	Order_ID       primitive.ObjectID `json:"order_id" bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	Order_Address  *string            `json:"order_address" bson:"order_address"`
	Order_Contact  *string            `json:"order_contact" bson:"order_contact"`
	Order_Status   *string            `json:"order_status" bson:"order_status"`
	Order_Total    *float64           `json:"order_total" bson:"order_total"`
	Order_Date     time.Time          `json:"order_date" bson:"order_date"`
	Price          *float64           `json:"price" bson:"price"`
	Discount       *float64           `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Address struct {
	Address_ID      primitive.ObjectID `json:"address_id" bson:"_id"`
	Address_Street  *string            `json:"address_street" bson:"address_street"`
	Address_City    *string            `json:"address_city" bson:"address_city"`
	Address_State   *string            `json:"address_state" bson:"address_state"`
	Address_Zip     *string            `json:"address_zip" bson:"address_zip"`
	Address_Country *string            `json:"address_country" bson:"address_country"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
