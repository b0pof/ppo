package user

import "git.iu7.bmstu.ru/kia22u475/ppo/internal/model"

type userRow struct {
	ID       int64   `db:"id"`
	Name     *string `db:"name"`
	Login    string  `db:"login"`
	Role     *string `db:"role"`
	Phone    *string `db:"phone"`
	Password string  `db:"password"`
}

func convertUser(row userRow) model.User {
	var name string
	if row.Name != nil {
		name = *row.Name
	}

	var role string
	if row.Role != nil {
		role = *row.Role
	}

	var phone string
	if row.Phone != nil {
		phone = *row.Phone
	}

	return model.User{
		ID:       row.ID,
		Name:     name,
		Login:    row.Login,
		Role:     role,
		Phone:    phone,
		Password: row.Password,
	}
}

type userMetaRow struct {
	Name            string `db:"name"`
	CartItemsAmount int    `db:"cart_items_amount"`
}

func convertUserMeta(row userMetaRow) model.UserMetaInfo {
	return model.UserMetaInfo{
		Name:            row.Name,
		CartItemsAmount: row.CartItemsAmount,
	}
}
