package handler

import (
	authHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/auth"
	cartHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/cart"
	categoryHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/category"
	itemHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/item"
	orderHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/order"
	reviewHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/review"
	sellerHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/seller"
	userHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/user"
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
