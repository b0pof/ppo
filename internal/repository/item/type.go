package item

import "git.iu7.bmstu.ru/kia22u475/ppo/internal/model"

type itemRow struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	SellerID    int64  `db:"seller_id"`
	SellerName  string `db:"seller_name"`
	Price       int    `db:"price"`
	Description string `db:"description"`
	ImgSrc      string `db:"imgsrc"`
}

func convertItem(row itemRow) model.Item {
	return model.Item{
		ID:   row.ID,
		Name: row.Name,
		Seller: model.Seller{
			ID:   row.SellerID,
			Name: row.SellerName,
		},
		Price:       row.Price,
		Description: row.Description,
		ImgSrc:      row.ImgSrc,
	}
}

func convertItems(rows []itemRow) []model.Item {
	items := make([]model.Item, 0, len(rows))

	for _, row := range rows {
		items = append(items, convertItem(row))
	}

	return items
}

type itemExtendedRow struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	SellerID    int64   `db:"seller_id"`
	Rating      *int    `db:"rating"`
	SellerName  *string `db:"seller_name"`
	Price       int     `db:"price"`
	Description string  `db:"description"`
	ImgSrc      string  `db:"imgsrc"`
	Amount      int     `db:"amount"`
}

func convertItemExtended(row itemExtendedRow) model.ItemExtended {
	sellerName := ""
	if row.SellerName != nil {
		sellerName = *row.SellerName
	}

	rating := 0
	if row.Rating != nil {
		rating = *row.Rating
	}

	return model.ItemExtended{
		Item: model.Item{
			ID:   row.ID,
			Name: row.Name,
			Seller: model.Seller{
				ID:   row.SellerID,
				Name: sellerName,
			},
			Rating:      float64(rating) / float64(100),
			Description: row.Description,
			Price:       row.Price,
			ImgSrc:      row.ImgSrc,
		},
		SellerName: sellerName,
		Amount:     row.Amount,
	}
}

func convertItemsExtended(rows []itemExtendedRow) []model.ItemExtended {
	items := make([]model.ItemExtended, 0, len(rows))

	for _, row := range rows {
		items = append(items, convertItemExtended(row))
	}

	return items
}

type orderItemRow struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	Price  int    `db:"price"`
	Count  int    `db:"count"`
	ImgSrc string `db:"imgsrc"`
}

func convertOrderItem(row orderItemRow) model.OrderItemInfo {
	return model.OrderItemInfo{
		ID:          row.ID,
		ProductName: row.Name,
		Price:       row.Price,
		Count:       row.Count,
		ImgSrc:      row.ImgSrc,
	}
}

func convertOrderItems(rows []orderItemRow) []model.OrderItemInfo {
	items := make([]model.OrderItemInfo, 0, len(rows))

	for _, row := range rows {
		items = append(items, convertOrderItem(row))
	}

	return items
}
