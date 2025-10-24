package controller

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.iu7.bmstu.ru/kia22u475/ppo/tests/db/client"
)

type Controller struct {
	t          provider.T
	db         *sqlx.DB
	redis      *redis.Client
	httpClient *http.Client
}

func NewDatabase(t provider.T) *sqlx.DB {
	return client.DbConnect()
}

func NewController(t provider.T) *Controller {
	redisClient := client.RedisConnect()

	client.CleanUpRedis(t, redisClient)

	db := client.DbConnect()

	err := client.CleanUpDB(t, db)
	require.NoError(t, err, "db cleanup failed")

	err = client.PrepareDB(t, db)
	require.NoError(t, err, "db preparation failed")

	return &Controller{
		t:          t,
		db:         db,
		redis:      redisClient,
		httpClient: &http.Client{},
	}
}

func (c *Controller) GetDB() *sqlx.DB {
	return c.db
}

func (c *Controller) GetRedis() *redis.Client {
	return c.redis
}

func (c *Controller) GetTesting() provider.T {
	return c.t
}

func (c *Controller) GetHttp() *http.Client {
	return c.httpClient
}

func (c *Controller) Finish() {
	assert.NoError(c.t, c.db.Close())
}
