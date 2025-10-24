package model

type Item struct {
	ID          int64
	Name        string
	Description string
	Rating      float64
	Price       int
	Seller      Seller
	ImgSrc      string
}

type ItemExtended struct {
	Item
	SellerName string
	Amount     int
}

type OrderItemInfo struct {
	ID          int64
	ProductName string
	Price       int
	Count       int
	ImgSrc      string
}

type CartContent struct {
	TotalPrice int
	TotalCount int
	Items      []CartItem
}

type CartItem struct {
	ID     int64
	Name   string
	Price  int
	Count  int
	ImgSrc string
}
