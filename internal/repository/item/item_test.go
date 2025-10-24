package item_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/item"
)

type ItemRepository struct {
	suite.Suite
}

func TestItemRepository(t *testing.T) {
	suite.RunSuite(t, new(ItemRepository))
}

func (s *ItemRepository) TestRepository_GetByID(t provider.T) {
	t.Parallel()

	testSeller := model.User{
		ID:   2,
		Name: "testSeller",
	}

	testItem := model.Item{
		ID:   1,
		Name: "test",
		Seller: model.Seller{
			ID:   testSeller.ID,
			Name: testSeller.Name,
		},
		Rating:      4,
		Price:       153,
		Description: "description",
		ImgSrc:      "https://test.com",
	}

	testCartItem := cartItem{
		userID: 999,
		itemID: 1,
		count:  15,
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) (itemID int64, userID int64)
		expectations func(a provider.Asserts, got model.ItemExtended, err error)
	}{
		{
			name: "success: item exists and in cart",
			prepare: func(db *sqlx.DB) (int64, int64) {
				insertItem(db, testItem)
				insertUser(db, testSeller)
				insertCartItem(db, testCartItem)

				return testItem.ID, 999
			},
			expectations: func(a provider.Asserts, got model.ItemExtended, err error) {
				a.NoError(err)

				a.Equal(model.ItemExtended{
					Item: model.Item{
						ID:          testItem.ID,
						Name:        testItem.Name,
						Description: testItem.Description,
						Rating:      testItem.Rating / float64(100),
						Price:       testItem.Price,
						Seller:      testItem.Seller,
						ImgSrc:      testItem.ImgSrc,
					},
					SellerName: testSeller.Name,
					Amount:     testCartItem.count,
				}, got)
			},
		},
		{
			name: "not found",
			prepare: func(db *sqlx.DB) (int64, int64) {
				return 99, 1
			},
			expectations: func(a provider.Asserts, got model.ItemExtended, err error) {
				a.Error(err)
				a.ErrorIs(err, model.ErrNotFound)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			itemID, userID := tc.prepare(db)

			repo := New(db)

			got, err := repo.GetByID(context.Background(), itemID, userID)

			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *ItemRepository) TestRepository_GetAllItems(t provider.T) {
	t.Parallel()

	testSeller := model.User{ID: 2, Name: "testSeller"}
	testItem := model.Item{
		ID: 1, Name: "test", Seller: model.Seller{ID: testSeller.ID, Name: testSeller.Name},
		Rating: 250, Price: 100, Description: "desc", ImgSrc: "https://x",
	}
	testCartItem := cartItem{userID: 999, itemID: 1, count: 3}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got []model.ItemExtended, err error)
	}{
		{
			name: "success: one item in cart",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testSeller)
				insertItem(db, testItem)
				insertCartItem(db, testCartItem)
				return testCartItem.userID
			},
			expectations: func(a provider.Asserts, got []model.ItemExtended, err error) {
				a.NoError(err)
				a.Len(got, 1)
				a.Equal(int64(1), got[0].ID)
				a.Equal(3, got[0].Amount)
				a.Equal("testSeller", got[0].SellerName)
			},
		},
		{
			name: "empty: no items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testSeller)
				return 999
			},
			expectations: func(a provider.Asserts, got []model.ItemExtended, err error) {
				a.NoError(err)
				a.Len(got, 0)
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

			got, err := repo.GetAllItems(context.Background(), userID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *ItemRepository) TestRepository_GetItemsBySellerID(t provider.T) {
	t.Parallel()

	testSeller := model.User{ID: 5, Name: "sellerX"}
	testItem := model.Item{
		ID: 10, Name: "itemX", Seller: model.Seller{ID: 5, Name: "sellerX"},
		Rating: 100, Price: 55, Description: "d", ImgSrc: "https://x",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got []model.Item, err error)
	}{
		{
			name: "success: has items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testSeller)
				insertItem(db, testItem)
				return 5
			},
			expectations: func(a provider.Asserts, got []model.Item, err error) {
				a.NoError(err)
				a.Len(got, 1)
				a.Equal("itemX", got[0].Name)
			},
		},
		{
			name: "empty: no items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testSeller)
				return 123
			},
			expectations: func(a provider.Asserts, got []model.Item, err error) {
				a.NoError(err)
				a.Len(got, 0)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			sellerID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetItemsBySellerID(context.Background(), sellerID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *ItemRepository) TestRepository_GetItemsByOrderID(t provider.T) {
	t.Parallel()

	testSeller := model.User{ID: 7, Name: "s"}
	testItem := model.Item{
		ID: 77, Name: "X", Seller: model.Seller{ID: 7, Name: "s"},
		Rating: 0, Price: 42, Description: "", ImgSrc: "https://x",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got []model.OrderItemInfo, err error)
	}{
		{
			name: "success: one item in order",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testSeller)
				insertItem(db, testItem)
				db.MustExec(`insert into order_item(order_id,item_id,count) values (1,77,2)`)
				return 1
			},
			expectations: func(a provider.Asserts, got []model.OrderItemInfo, err error) {
				a.NoError(err)
				a.Len(got, 1)
				a.Equal(int64(77), got[0].ID)
				a.Equal(2, got[0].Count)
			},
		},
		{
			name: "empty",
			prepare: func(db *sqlx.DB) int64 {
				return 99
			},
			expectations: func(a provider.Asserts, got []model.OrderItemInfo, err error) {
				a.NoError(err)
				a.Len(got, 0)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			orderID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetItemsByOrderID(context.Background(), orderID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *ItemRepository) TestRepository_GetItemsByCategoryID(t provider.T) {
	t.Parallel()

	testSeller := model.User{ID: 9, Name: "s"}
	testItem := model.Item{
		ID: 99, Name: "itemC", Seller: model.Seller{ID: 9, Name: "s"},
		Rating: 10, Price: 33, Description: "d", ImgSrc: "https://x",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) (catID int64, userID int64)
		expectations func(a provider.Asserts, got []model.ItemExtended, err error)
	}{
		{
			name: "success: one item in category",
			prepare: func(db *sqlx.DB) (int64, int64) {
				insertUser(db, testSeller)
				insertItem(db, testItem)
				db.MustExec(`insert into item_category(item_id,category_id) values (99,5)`)
				return 5, 111
			},
			expectations: func(a provider.Asserts, got []model.ItemExtended, err error) {
				a.NoError(err)
				a.Len(got, 1)
				a.Equal(int64(99), got[0].ID)
			},
		},
		{
			name: "empty",
			prepare: func(db *sqlx.DB) (int64, int64) {
				return 5, 111
			},
			expectations: func(a provider.Asserts, got []model.ItemExtended, err error) {
				a.NoError(err)
				a.Len(got, 0)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			catID, userID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetItemsByCategoryID(context.Background(), catID, userID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func insertUser(db *sqlx.DB, user model.User) int64 {
	var newID int64

	_ = db.Get(
		&newID,
		`insert into "user" (id, role, name, login, phone, password) values ($1, $2, $3, $4, $5, $6) returning id`,
		user.ID,
		user.Role,
		user.Name,
		user.Login,
		user.Phone,
		user.Password,
	)

	return newID
}

func insertItem(db *sqlx.DB, item model.Item) int64 {
	var newID int64

	_ = db.Get(
		&newID,
		`insert into item (id, name, seller_id, rating, price, description, imgsrc) values ($1, $2, $3, $4, $5, $6, $7) returning id`,
		item.ID,
		item.Name,
		item.Seller.ID,
		item.Rating,
		item.Price,
		item.Description,
		item.ImgSrc,
	)

	return newID
}

type cartItem struct {
	userID int64
	itemID int64
	count  int
}

func insertCartItem(db *sqlx.DB, item cartItem) {
	_ = db.MustExec(
		`insert into cart_item(user_id, item_id, count) values ($1, $2, $3)`,
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
			id bigserial,
			role text,
			name text not null,
			login text not null unique,
			phone text,
			password text not null
		);
		create table if not exists item (
			id bigserial,
			name text not null,
			seller_id bigint,
			price int,
			rating int,
			description text default '',
			imgsrc text not null
		);
		create table if not exists cart_item (
			user_id bigint,
			item_id bigint,
			count int not null default 1,
			constraint unique_user_id_item_id unique (user_id, item_id)
		);
		create table if not exists order_item (
			order_id bigint,
			item_id bigint,
			count int not null
		);
		create table if not exists item_category (
			item_id bigint,
			category_id bigint
		);
	`
	db.MustExec(schema)
	return db
}
