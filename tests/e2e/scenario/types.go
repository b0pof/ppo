package scenario

// create_order_flow

type loginResponse struct {
	SessionID string `json:"sessionID"`
	UserID    int64  `json:"userID"`
}

type createOrderResponse struct {
	OrderID int64 `json:"orderID"`
}

type cartItemsResponse struct {
	TotalPrice int        `json:"totalPrice"`
	TotalCount int        `json:"totalCount"`
	Items      []cartItem `json:"items"`
}

type cartItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Count int    `json:"count"`
}

type orderResponse struct {
	ID         int64  `json:"id"`
	BuyerID    int64  `json:"buyerId"`
	ItemsCount int    `json:"itemsCount"`
	Status     string `json:"status"`
	Sum        int    `json:"sum"`
	Items      []struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
		Count int    `json:"count"`
	} `json:"items"`
}
