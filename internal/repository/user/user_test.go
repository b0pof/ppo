package user_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/user"
)

type UserRepository struct {
	suite.Suite
}

func TestUserRepository(t *testing.T) {
	suite.RunSuite(t, new(UserRepository))
}

func (s *UserRepository) TestRepository_GetByID(t provider.T) {
	t.Parallel()

	testUser := model.User{
		ID:       1,
		Name:     "Test User",
		Login:    "testuser",
		Role:     "customer",
		Phone:    "+1234567890",
		Password: "hashedpassword",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got model.User, err error)
	}{
		{
			name: "success: user exists",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				return testUser.ID
			},
			expectations: func(a provider.Asserts, got model.User, err error) {
				a.NoError(err)
				a.Equal(testUser, got)
			},
		},
		{
			name: "not found",
			prepare: func(db *sqlx.DB) int64 {
				return 999
			},
			expectations: func(a provider.Asserts, got model.User, err error) {
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

			userID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetByID(context.Background(), userID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *UserRepository) TestRepository_GetByLogin(t provider.T) {
	t.Parallel()

	testUser := model.User{
		ID:       1,
		Name:     "Test User",
		Login:    "testuser",
		Role:     "customer",
		Phone:    "+1234567890",
		Password: "hashedpassword",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) string
		expectations func(a provider.Asserts, got model.User, err error)
	}{
		{
			name: "success: user exists",
			prepare: func(db *sqlx.DB) string {
				insertUser(db, testUser)
				return testUser.Login
			},
			expectations: func(a provider.Asserts, got model.User, err error) {
				a.NoError(err)
				a.Equal(testUser, got)
			},
		},
		{
			name: "not found",
			prepare: func(db *sqlx.DB) string {
				return "nonexistent"
			},
			expectations: func(a provider.Asserts, got model.User, err error) {
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

			login := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetByLogin(context.Background(), login)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *UserRepository) TestRepository_GetUserMetaByID(t provider.T) {
	t.Parallel()

	testUser := model.User{
		ID:       1,
		Name:     "Test User",
		Login:    "testuser",
		Role:     "customer",
		Phone:    "+1234567890",
		Password: "hashedpassword",
	}

	testItem := model.Item{
		ID:   1,
		Name: "Test Item",
		Seller: model.Seller{
			ID: 2, // другой продавец
		},
		Price:  100,
		ImgSrc: "test.jpg",
	}

	testSeller := model.User{
		ID:       2,
		Name:     "Test Seller",
		Login:    "testseller",
		Role:     "seller",
		Phone:    "+0987654321",
		Password: "sellerpass",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got model.UserMetaInfo, err error)
	}{
		{
			name: "success: user with cart items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				insertUser(db, testSeller)
				insertItem(db, testItem)
				insertCartItem(db, cartItem{userID: testUser.ID, itemID: testItem.ID, count: 3})
				return testUser.ID
			},
			expectations: func(a provider.Asserts, got model.UserMetaInfo, err error) {
				a.NoError(err)
				a.Equal(testUser.Name, got.Name)
				a.Equal(1, got.CartItemsAmount) // количество уникальных товаров в корзине
			},
		},
		{
			name: "success: user without cart items",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				return testUser.ID
			},
			expectations: func(a provider.Asserts, got model.UserMetaInfo, err error) {
				a.NoError(err)
				a.Equal(testUser.Name, got.Name)
				a.Equal(0, got.CartItemsAmount)
			},
		},
		{
			name: "not found",
			prepare: func(db *sqlx.DB) int64 {
				return 999
			},
			expectations: func(a provider.Asserts, got model.UserMetaInfo, err error) {
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

			userID := tc.prepare(db)
			repo := New(db)

			got, err := repo.GetUserMetaByID(context.Background(), userID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *UserRepository) TestRepository_GetUserLoginByID(t provider.T) {
	t.Parallel()

	testUser := model.User{
		ID:       1,
		Name:     "Test User",
		Login:    "testuser",
		Role:     "customer",
		Phone:    "+1234567890",
		Password: "hashedpassword",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got string, err error)
	}{
		{
			name: "success: user exists",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				return testUser.ID
			},
			expectations: func(a provider.Asserts, got string, err error) {
				a.NoError(err)
				a.Equal(testUser.Login, got)
			},
		},
		{
			name: "not found",
			prepare: func(db *sqlx.DB) int64 {
				return 999
			},
			expectations: func(a provider.Asserts, got string, err error) {
				a.Error(err)
				a.Contains(err.Error(), "user repository.GetUserLoginByID")
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

			got, err := repo.GetUserLoginByID(context.Background(), userID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *UserRepository) TestRepository_GetRoleByID(t provider.T) {
	t.Parallel()

	testUser := model.User{
		ID:       1,
		Name:     "Test User",
		Login:    "testuser",
		Role:     "customer",
		Phone:    "+1234567890",
		Password: "hashedpassword",
	}

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB) int64
		expectations func(a provider.Asserts, got string, err error)
	}{
		{
			name: "success: user exists",
			prepare: func(db *sqlx.DB) int64 {
				insertUser(db, testUser)
				return testUser.ID
			},
			expectations: func(a provider.Asserts, got string, err error) {
				a.NoError(err)
				a.Equal(testUser.Role, got)
			},
		},
		{
			name: "not found",
			prepare: func(db *sqlx.DB) int64 {
				return 999
			},
			expectations: func(a provider.Asserts, got string, err error) {
				a.Error(err)
				a.Contains(err.Error(), "user repository.GetRoleByID")
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

			got, err := repo.GetRoleByID(context.Background(), userID)
			tc.expectations(ctx.Assert(), got, err)
		})
	}
}

func (s *UserRepository) TestRepository_Create(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		prepare      func(db *sqlx.DB)
		login        string
		password     string
		role         string
		expectations func(a provider.Asserts, id int64, err error, db *sqlx.DB)
	}{
		{
			name:     "success: new user",
			prepare:  func(db *sqlx.DB) {},
			login:    "newuser",
			password: "password123",
			role:     "customer",
			expectations: func(a provider.Asserts, id int64, err error, db *sqlx.DB) {
				a.Error(err)
			},
		},
		{
			name: "already exists",
			prepare: func(db *sqlx.DB) {
				insertUser(db, model.User{
					ID:       1,
					Name:     "Existing User",
					Login:    "existinguser",
					Role:     "customer",
					Password: "oldpassword",
				})
			},
			login:    "existinguser",
			password: "newpassword",
			role:     "seller",
			expectations: func(a provider.Asserts, id int64, err error, db *sqlx.DB) {
				a.Error(err)
				a.Equal(int64(0), id)

				var password string
				err = db.Get(&password, `SELECT password FROM "user" WHERE login = ?`, "existinguser")
				a.NoError(err)
				a.Equal("oldpassword", password)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			db := mustInitDB()
			defer db.Close()

			tc.prepare(db)
			repo := New(db)

			id, err := repo.Create(context.Background(), tc.login, tc.password, tc.role)
			tc.expectations(ctx.Assert(), id, err, db)
		})
	}
}

func insertUser(db *sqlx.DB, user model.User) int64 {
	res, err := db.Exec(
		`INSERT INTO "user" (id, role, name, login, phone, password) VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Role, user.Name, user.Login, user.Phone, user.Password)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}

func insertItem(db *sqlx.DB, item model.Item) int64 {
	res, err := db.Exec(
		`INSERT INTO item (id, name, seller_id, price, imgsrc) VALUES (?, ?, ?, ?, ?)`,
		item.ID, item.Name, item.Seller.ID, item.Price, item.ImgSrc)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}

type cartItem struct {
	userID int64
	itemID int64
	count  int
}

func insertCartItem(db *sqlx.DB, item cartItem) {
	_, err := db.Exec(
		`INSERT INTO cart_item (user_id, item_id, count) VALUES (?, ?, ?)`,
		item.userID, item.itemID, item.count)
	if err != nil {
		panic(err)
	}
}

func mustInitDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	schema := `
		create table if not exists item (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
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
		
		create table if not exists "user" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role text,
			name text not null,
			login text not null unique,
			phone text,
			password text not null
		);
		create table if not exists order_item (
			order_id bigint,
			item_id bigint,
			count int
		);
		create table if not exists item_category (
			item_id bigint,
			category_id bigint
		);
	`
	db.MustExec(schema)
	return db
}
