package cli

import (
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

const (
	userLayout = `
	ID: %d
	Логин: %s
	Имя: %s
	Роль: %s
	Телефон: %s
	`

	itemExtendedLayout = `
	ID: %d
	Название: %s
	Описание: %s
	Цена: %d
	Продавец: %s
	Картинка: %s
	В корзине: %d
	`

	itemLayout = `
	ID: %d
	Название: %s
	Описание: %s
	Цена: %d
	Продавец: %s
	Картинка: %s
	`

	orderLayout = `
	ID: %d
	Сумма: %d
	ID покупателя: %d
	Кол-во товаров: %d
	Статус: %s
	Время создания: %s
	`
)

func formatUser(u model.User) string {
	return fmt.Sprintf(
		userLayout,
		u.ID,
		u.Login,
		u.Name,
		u.Role,
		u.Phone,
	)
}

func formatItemExtended(item model.ItemExtended) string {
	return fmt.Sprintf(
		itemExtendedLayout,
		item.ID,
		item.Name,
		item.Description,
		item.Price,
		item.Seller.Name,
		item.ImgSrc,
		item.Amount,
	)
}

func formatItem(item model.Item) string {
	return fmt.Sprintf(
		itemLayout,
		item.ID,
		item.Name,
		item.Description,
		item.Price,
		item.Seller.Name,
		item.ImgSrc,
	)
}

func formatOrder(o model.Order) string {
	return fmt.Sprintf(
		orderLayout,
		o.ID,
		o.Sum,
		o.BuyerID,
		o.ItemsCount,
		o.Status,
		o.CreatedAt,
	)
}
