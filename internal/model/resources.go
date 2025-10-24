package model

import (
	"fmt"
	"strings"
)

type Path string

const (
	LoginPath     Path = "POST:/api/1/auth"
	SignupPath    Path = "POST:/api/1/users"
	LogoutPath    Path = "DELETE:/api/1/auth"
	AuthCheckPath Path = "GET:/api/1/auth"

	CartItemAddPath    Path = "POST:/api/1/users/{id}/cart/items"
	CartItemDeletePath Path = "DELETE:/api/1/users/{userId}/cart/items/{itemId}"
	CartClearPath      Path = "DELETE:/api/1/users/{id}/cart/items"
	CartItemsPath      Path = "GET:/api/1/users/{id}/cart/items"

	ItemCreatePath Path = "POST:/api/1/items"
	ItemFetchPath  Path = "GET:/api/1/items/{id}"
	ItemsFetchPath Path = "GET:/api/1/items"
	ItemDeletePath Path = "DELETE:/api/1/items/{id}"
	ItemUpdatePath Path = "UPDATE:/api/1/items/{id}"

	OrderCreatePath       Path = "POST:/api/1/users/{id}/orders"
	OrderPath             Path = "GET:/api/1/users/{userId}/orders/{orderId}"
	OrdersPath            Path = "GET:/api/1/users/{id}/orders"
	OrderStatusUpdatePath Path = "PATCH:/api/1/users/{userId}/orders/{orderId}"

	SellerItemsFetchPath Path = "GET:/api/1/seller/{id}/items"

	UserFetchPath          Path = "GET:/api/1/users/{id}"
	UserMetaFetchPath      Path = "GET:/api/1/users/{id}/meta"
	UserUpdatePath         Path = "PUT:/api/1/users/{id}"
	UserPasswordUpdatePath Path = "PATCH:/api/1/users/{id}/password"

	CategoryFetchPath      Path = "GET:/api/1/categories/{id}"
	CategoryItemsFetchPath Path = "GET:/api/1/categories/{id}/items"

	ReviewsFetchPath Path = "GET:/api/1/items/{id}/reviews"
	ReviewCreatePath Path = "POST:/api/1/items/{id}/reviews"
)

var Resources = map[Path]Permissions{
	// Auth
	SignupPath: {
		{
			role: RoleGuest,
		},
	},
	LoginPath: {
		{
			role: RoleGuest,
		},
	},
	LogoutPath: {
		{
			role: RoleSeller,
		},
		{
			role: RoleBuyer,
		},
	},
	AuthCheckPath: {
		{
			role: RoleGuest,
		},
		{
			role: RoleSeller,
		},
		{
			role: RoleBuyer,
		},
	},

	// Cart
	CartItemAddPath: {
		{
			role: RoleBuyer,
		},
	},
	CartItemDeletePath: {
		{
			role: RoleBuyer,
		},
	},
	CartClearPath: {
		{
			role: RoleBuyer,
		},
	},
	CartItemsPath: {
		{
			role: RoleBuyer,
		},
	},

	// Item
	ItemCreatePath: {
		{
			role: RoleSeller,
		},
	},
	ItemFetchPath: {
		{
			role: RoleGuest,
		},
		{
			role: RoleSeller,
		},
		{
			role: RoleBuyer,
		},
	},
	ItemsFetchPath: {
		{
			role: RoleGuest,
		},
		{
			role: RoleSeller,
		},
		{
			role: RoleBuyer,
		},
	},
	ItemDeletePath: {
		{
			role: RoleSeller,
		},
	},
	ItemUpdatePath: {
		{
			role: RoleSeller,
		},
	},

	// Order
	OrderCreatePath: {
		{
			role: RoleBuyer,
		},
	},
	OrderPath: {
		{
			role: RoleBuyer,
		},
	},
	OrdersPath: {
		{
			role: RoleBuyer,
		},
	},
	OrderStatusUpdatePath: {
		{
			role: RoleGuest, // Postman
		},
		{
			role: RoleSeller,
		},
	},

	// User
	SellerItemsFetchPath: {
		{
			role: RoleSeller,
		},
	},
	UserFetchPath: {
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
	},
	UserMetaFetchPath: {
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
	},
	UserUpdatePath: {
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
	},
	UserPasswordUpdatePath: {
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
	},

	// Category
	CategoryFetchPath: {
		{
			role: RoleGuest,
		},
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
	},
	CategoryItemsFetchPath: {
		{
			role: RoleGuest,
		},
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
	},

	// Review
	ReviewCreatePath: {
		{
			role: RoleBuyer,
		},
	},
	ReviewsFetchPath: {
		{
			role: RoleBuyer,
		},
		{
			role: RoleSeller,
		},
		{
			role: RoleGuest,
		},
	},
}

func NewPath(path, method string) Path {
	return Path(fmt.Sprintf("%s:%s", method, path))
}

func (p Path) Method() string {
	if len(string(p)) == 0 {
		return ""
	}

	parts := strings.Split(string(p), ":")
	if len(parts) != 2 {
		return ""
	}

	return parts[0]
}

func (p Path) Url() string {
	if len(string(p)) == 0 {
		return ""
	}

	parts := strings.Split(string(p), ":")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}
