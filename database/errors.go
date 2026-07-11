package database

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserIdIsNotValid = errors.New("user id is not valid")
	ErrUpdateTokens     = errors.New("failed to update tokens")
	ErrUpdateUser       = errors.New("failed to update user")

	ErrCantFindProduct = errors.New("can't find product")
	ErrAddToCart       = errors.New("cannot add product to cart")
	ErrRemoveFromCart  = errors.New("cannot remove item from cart")
	ErrGetCart         = errors.New("unable to get cart")
	ErrCheckoutCart    = errors.New("cannot complete purchase")

	ErrTooManyAddresses = errors.New("cannot add more than 2 addresses")
	ErrAddAddress       = errors.New("failed to add address")
	ErrUpdateAddress    = errors.New("failed to update address")
	ErrDeleteAddress    = errors.New("failed to delete addresses")

	ErrCreateUser     = errors.New("user was not created")
	ErrCreateProduct  = errors.New("product was not created")
	ErrListProducts   = errors.New("unable to fetch products")
	ErrSearchProducts = errors.New("unable to search products")
)
