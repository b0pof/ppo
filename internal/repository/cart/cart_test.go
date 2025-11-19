package cart_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	. "github.com/b0pof/ppo/internal/repository/cart"
)

type CartRepositorySuite struct {
	suite.Suite
}

func TestCartRepository(t *testing.T) {
	suite.RunSuite(t, new(CartRepositorySuite))
}

func (s *CartRepositorySuite) TestRepository_GetCartItemsAmount(t provider.T) {
	t.Parallel()

	testUser := model.User{ID: 1, Name: "testUser"}
	testSeller := model.User{ID: 2, Name: "testSeller"}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got int, err error)
	}{
		{
			name: "success: user has cart items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)

				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})
				insertItem(db, model.Item{ID: 2, Name: "item2", Seller: model.Seller{ID: testSeller.ID}, Price: 200})

				insertCartItem(db, cartItem{userID: 1, itemID: 1, count: 2})
				insertCartItem(db, cartItem{userID: 1, itemID: 2, count: 1})

				return 1
			},
			expectations: func(a provider.Asserts, got int, err error) {
				a.NoError(err)
				a.Equal(2, got)
			},
		},
		{
			name: "success: user has no cart items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)
				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})

				insertCartItem(db, cartItem{userID: 999, itemID: 1, count: 5})

				return 1
			},
			expectations: func(a provider.Asserts, got int, err error) {
				a.NoError(err)
				a.Equal(0, got)
			},
		},
		{
			name: "success: user doesn't exist",
			prepare: func(db *sqlx.DB) int64 {
				return 9999
			},
			expectations: func(a provider.Asserts, got int, err error) {
				a.NoError(err)
				a.Equal(0, got)
			},
		},
		{
			name: "success: multiple items with different quantities",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)

				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})
				insertItem(db, model.Item{ID: 2, Name: "item2", Seller: model.Seller{ID: testSeller.ID}, Price: 200})
				insertItem(db, model.Item{ID: 3, Name: "item3", Seller: model.Seller{ID: testSeller.ID}, Price: 300})

				insertCartItem(db, cartItem{userID: 1, itemID: 1, count: 5})
				insertCartItem(db, cartItem{userID: 1, itemID: 2, count: 3})
				insertCartItem(db, cartItem{userID: 1, itemID: 3, count: 1})

				return 1
			},
			expectations: func(a provider.Asserts, got int, err error) {
				a.NoError(err)
				a.Equal(3, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			userID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetCartItemsAmount(context.Background(), userID)

			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *CartRepositorySuite) TestRepository_GetCartContentByUserID(t provider.T) {
	t.Parallel()

	testUser := model.User{ID: 1, Name: "testUser"}
	testSeller := model.User{ID: 2, Name: "testSeller"}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got model.CartContent, err error)
	}{
		{
			name: "success: user has items in cart",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)

				insertItem(db, model.Item{
					ID:     1,
					Name:   "Laptop",
					Seller: model.Seller{ID: testSeller.ID},
					Price:  1000,
					ImgSrc: "laptop.jpg",
				})
				insertItem(db, model.Item{
					ID:     2,
					Name:   "Mouse",
					Seller: model.Seller{ID: testSeller.ID},
					Price:  50,
					ImgSrc: "mouse.jpg",
				})

				insertCartItem(db, cartItem{userID: 1, itemID: 1, count: 1})
				insertCartItem(db, cartItem{userID: 1, itemID: 2, count: 2})

				return 1
			},
			expectations: func(a provider.Asserts, got model.CartContent, err error) {
				a.NoError(err)
				a.Equal(1100, got.TotalPrice)
				a.Equal(3, got.TotalCount)
				a.Len(got.Items, 2)

				a.Equal(int64(1), got.Items[0].ID)
				a.Equal("Laptop", got.Items[0].Name)
				a.Equal(1000, got.Items[0].Price)
				a.Equal(1, got.Items[0].Count)
				a.Equal("laptop.jpg", got.Items[0].ImgSrc)

				a.Equal(int64(2), got.Items[1].ID)
				a.Equal("Mouse", got.Items[1].Name)
				a.Equal(50, got.Items[1].Price)
				a.Equal(2, got.Items[1].Count)
				a.Equal("mouse.jpg", got.Items[1].ImgSrc)
			},
		},
		{
			name: "success: empty cart",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)
				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})

				insertCartItem(db, cartItem{userID: 999, itemID: 1, count: 5})

				return 1
			},
			expectations: func(a provider.Asserts, got model.CartContent, err error) {
				a.Error(err)
				a.ErrorIs(err, model.ErrNotFound)
				a.Equal(model.CartContent{}, got)
			},
		},
		{
			name: "success: single item with large quantity",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)

				insertItem(db, model.Item{
					ID:     1,
					Name:   "Book",
					Seller: model.Seller{ID: testSeller.ID},
					Price:  20,
					ImgSrc: "book.jpg",
				})

				insertCartItem(db, cartItem{userID: 1, itemID: 1, count: 10})

				return 1
			},
			expectations: func(a provider.Asserts, got model.CartContent, err error) {
				a.NoError(err)
				a.Equal(200, got.TotalPrice)
				a.Equal(10, got.TotalCount)
				a.Len(got.Items, 1)
				a.Equal(int64(1), got.Items[0].ID)
				a.Equal(10, got.Items[0].Count)
			},
		},
		{
			name: "success: multiple items with same price",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)

				insertItem(db, model.Item{ID: 1, Name: "ItemA", Seller: model.Seller{ID: testSeller.ID}, Price: 100, ImgSrc: "a.jpg"})
				insertItem(db, model.Item{ID: 2, Name: "ItemB", Seller: model.Seller{ID: testSeller.ID}, Price: 100, ImgSrc: "b.jpg"})
				insertItem(db, model.Item{ID: 3, Name: "ItemC", Seller: model.Seller{ID: testSeller.ID}, Price: 100, ImgSrc: "c.jpg"})

				insertCartItem(db, cartItem{userID: 1, itemID: 1, count: 1})
				insertCartItem(db, cartItem{userID: 1, itemID: 2, count: 1})
				insertCartItem(db, cartItem{userID: 1, itemID: 3, count: 1})

				return 1
			},
			expectations: func(a provider.Asserts, got model.CartContent, err error) {
				a.NoError(err)
				a.Equal(300, got.TotalPrice)
				a.Equal(3, got.TotalCount)
				a.Len(got.Items, 3)
			},
		},
		{
			name: "error: user doesn't exist",
			prepare: func(db *sqlx.DB) int64 {
				// Не создаем пользователя
				return 9999
			},
			expectations: func(a provider.Asserts, got model.CartContent, err error) {
				a.Error(err)
				a.ErrorIs(err, model.ErrNotFound)
				a.Equal(model.CartContent{}, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			userID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetCartContentByUserID(context.Background(), userID)

			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

type cartItem struct {
	userID int64
	itemID int64
	count  int
}

func insertUser(db *sqlx.DB, user model.User) {
	db.MustExec(
		`insert into "user" (id, role, name, login, phone, password) values ($1, $2, $3, $4, $5, $6)`,
		user.ID,
		user.Role,
		user.Name,
		user.Login,
		user.Phone,
		user.Password,
	)
}

func insertItem(db *sqlx.DB, item model.Item) {
	db.MustExec(
		`insert into item (id, name, seller_id, rating, price, description, imgsrc) values ($1, $2, $3, $4, $5, $6, $7)`,
		item.ID,
		item.Name,
		item.Seller.ID,
		item.Rating,
		item.Price,
		item.Description,
		item.ImgSrc,
	)
}

func insertCartItem(db *sqlx.DB, item cartItem) {
	db.MustExec(
		`insert into cart_item (user_id, item_id, count) values ($1, $2, $3)`,
		item.userID,
		item.itemID,
		item.count,
	)
}

func mustInitDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	schema := `
		create table if not exists "user" (
			id integer primary key,
			role text,
			name text not null,
			login text not null, -- unique,
			phone text,
			password text not null
		);
		create table if not exists item (
			id integer primary key,
			name text not null,
			seller_id integer,
			price integer,
			rating integer,
			description text default '',
			imgsrc text not null
		);
		create table if not exists cart_item (
			user_id integer,
			item_id integer,
			count integer not null default 1,
			primary key (user_id, item_id)
		);
	`

	db.MustExec(schema)
	return db
}
