package cart

import "github.com/b0pof/ppo/internal/model"

type cartItemRow struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	Price  int    `db:"price"`
	Count  int    `db:"count"`
	ImgSrc string `db:"imgsrc"`
}

func convertCartItem(row cartItemRow) model.CartItem {
	return model.CartItem{
		ID:     row.ID,
		Name:   row.Name,
		Price:  row.Price,
		Count:  row.Count,
		ImgSrc: row.ImgSrc,
	}
}

func convertCartItems(rows []cartItemRow) []model.CartItem {
	items := make([]model.CartItem, 0, len(rows))

	for _, row := range rows {
		items = append(items, convertCartItem(row))
	}

	return items
}
