package handler

import (
	authHandler "github.com/b0pof/ppo/internal/api/handler/auth"
	cartHandler "github.com/b0pof/ppo/internal/api/handler/cart"
	categoryHandler "github.com/b0pof/ppo/internal/api/handler/category"
	itemHandler "github.com/b0pof/ppo/internal/api/handler/item"
	orderHandler "github.com/b0pof/ppo/internal/api/handler/order"
	reviewHandler "github.com/b0pof/ppo/internal/api/handler/review"
	sellerHandler "github.com/b0pof/ppo/internal/api/handler/seller"
	userHandler "github.com/b0pof/ppo/internal/api/handler/user"
)

type Handler struct {
	*authHandler.Auth
	*cartHandler.Cart
	*categoryHandler.Category
	*itemHandler.Item
	*orderHandler.Order
	*reviewHandler.Review
	*sellerHandler.Seller
	*userHandler.User
}

func NewHandler(
	auth *authHandler.Auth,
	cart *cartHandler.Cart,
	category *categoryHandler.Category,
	item *itemHandler.Item,
	order *orderHandler.Order,
	review *reviewHandler.Review,
	seller *sellerHandler.Seller,
	user *userHandler.User,
) *Handler {
	return &Handler{
		Auth:     auth,
		Cart:     cart,
		Category: category,
		Item:     item,
		Order:    order,
		Review:   review,
		Seller:   seller,
		User:     user,
	}
}
