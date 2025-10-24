package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/order"
)

type OrderRepositorySuite struct {
	suite.Suite
}

func TestOrderRepository(t *testing.T) {
	suite.RunSuite(t, new(OrderRepositorySuite))
}

func (s *OrderRepositorySuite) TestRepository_GetByID(t provider.T) {
	t.Parallel()

	testBuyer := model.User{
		ID:   1,
		Name: "testBuyer",
	}

	testSeller := model.User{
		ID:   2,
		Name: "testSeller",
	}

	testItem := model.Item{
		ID:     1,
		Name:   "testItem",
		Seller: model.Seller{ID: testSeller.ID, Name: testSeller.Name},
		Price:  100,
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got model.Order, err error)
	}{
		{
			name: "success: order exists with items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testBuyer)
				insertUser(db, testSeller)
				insertItem(db, testItem)

				var orderID int64
				db.Get(&orderID, `insert into "order" (id, buyer_id, status, created_at) values (1, $1, 'created', $2) returning id`,
					testBuyer.ID, time.Now())

				db.MustExec(`insert into order_item (order_id, item_id, count) values (1, 1, 2)`)
				db.MustExec(`insert into order_item (order_id, item_id, count) values (1, 1, 1)`)

				return 1
			},
			expectations: func(a provider.Asserts, got model.Order, err error) {
				a.NoError(err)
				a.Equal(int64(1), got.ID)
				a.Equal(testBuyer.ID, got.BuyerID)
				a.Equal(model.OrderStatusCreated, got.Status)
				a.Equal(300, got.Sum)
				a.NotZero(got.CreatedAt)
			},
		},
		{
			name: "not found: order doesn't exist",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testBuyer)
				return 999
			},
			expectations: func(a provider.Asserts, got model.Order, err error) {
				a.Error(err)
				a.ErrorIs(err, model.ErrNotFound)
				a.Equal(model.Order{}, got)
			},
		},
		{
			name: "success: order with multiple items different prices",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testBuyer)
				insertUser(db, testSeller)

				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})
				insertItem(db, model.Item{ID: 2, Name: "item2", Seller: model.Seller{ID: testSeller.ID}, Price: 200})

				var orderID int64
				db.Get(&orderID, `insert into "order" (id, buyer_id, status, created_at) values (2, $1, 'cancelled', $2) returning id`,
					testBuyer.ID, time.Now())

				db.MustExec(`insert into order_item (order_id, item_id, count) values (2, 1, 3)`)
				db.MustExec(`insert into order_item (order_id, item_id, count) values (2, 2, 1)`)

				return 2
			},
			expectations: func(a provider.Asserts, got model.Order, err error) {
				a.NoError(err)
				a.Equal(int64(2), got.ID)
				a.Equal(500, got.Sum)
				a.Equal(model.OrderStatusCancelled, got.Status)
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

			got, err := repo.GetByID(context.Background(), orderID)

			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *OrderRepositorySuite) TestRepository_GetOrdersByUserID(t provider.T) {
	t.Parallel()

	testBuyer1 := model.User{ID: 1, Name: "buyer1"}
	testBuyer2 := model.User{ID: 2, Name: "buyer2"}
	testSeller := model.User{ID: 3, Name: "seller"}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got []model.Order, err error)
	}{
		{
			name: "success: user has multiple orders",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testBuyer1)
				insertUser(db, testBuyer2)
				insertUser(db, testSeller)

				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})
				insertItem(db, model.Item{ID: 2, Name: "item2", Seller: model.Seller{ID: testSeller.ID}, Price: 150})

				db.MustExec(`insert into "order" (id, buyer_id, status, created_at) values (1, 1, 'created', $1)`, time.Now())
				db.MustExec(`insert into "order" (id, buyer_id, status, created_at) values (2, 1, 'done', $1)`, time.Now().Add(-time.Hour))

				db.MustExec(`insert into "order" (id, buyer_id, status, created_at) values (3, 2, 'created', $1)`, time.Now())

				db.MustExec(`insert into order_item (order_id, item_id, count) values (1, 1, 2)`)
				db.MustExec(`insert into order_item (order_id, item_id, count) values (2, 2, 1)`)

				return 1
			},
			expectations: func(a provider.Asserts, got []model.Order, err error) {
				a.NoError(err)
				a.Len(got, 2)

				a.Equal(int64(1), got[0].ID)
				a.Equal(200, got[0].Sum)
				a.Equal(model.OrderStatusCreated, got[0].Status)

				a.Equal(int64(2), got[1].ID)
				a.Equal(150, got[1].Sum)
				a.Equal(model.OrderStatusDone, got[1].Status)
			},
		},
		{
			name: "success: user has no orders",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testBuyer1)
				insertUser(db, testSeller)
				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})

				db.MustExec(`insert into "order" (id, buyer_id, status, created_at) values (1, 999, 'pending', $1)`, time.Now())
				db.MustExec(`insert into order_item (order_id, item_id, count) values (1, 1, 1)`)

				return 1
			},
			expectations: func(a provider.Asserts, got []model.Order, err error) {
				a.ErrorIs(err, model.ErrNotFound)
				a.Len(got, 0)
			},
		},
		{
			name: "not found: user doesn't exist",
			prepare: func(db *sqlx.DB) int64 {
				return 9997
			},
			expectations: func(a provider.Asserts, got []model.Order, err error) {
				a.Error(err)
				a.ErrorIs(err, model.ErrNotFound)
				a.Nil(got)
			},
		},
		{
			name: "success: order with multiple items per order",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testBuyer1)
				insertUser(db, testSeller)

				insertItem(db, model.Item{ID: 1, Name: "item1", Seller: model.Seller{ID: testSeller.ID}, Price: 100})
				insertItem(db, model.Item{ID: 2, Name: "item2", Seller: model.Seller{ID: testSeller.ID}, Price: 200})

				db.MustExec(`insert into "order" (id, buyer_id, status, created_at) values (1, 1, 'pending', $1)`, time.Now())

				db.MustExec(`insert into order_item (order_id, item_id, count) values (1, 1, 2)`)
				db.MustExec(`insert into order_item (order_id, item_id, count) values (1, 2, 1)`)

				return 1
			},
			expectations: func(a provider.Asserts, got []model.Order, err error) {
				a.NoError(err)
				a.Len(got, 1)
				a.Equal(int64(1), got[0].ID)
				a.Equal(400, got[0].Sum)
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

			got, err := repo.GetOrdersByUserID(context.Background(), userID)

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
		create table if not exists "order" (
			id bigserial,
			buyer_id bigint,
			status text not null,
			created_at timestamp not null
		);
		create table if not exists order_item (
			order_id bigint,
			item_id bigint,
			count int not null
		);
	`
	db.MustExec(schema)
	return db
}
