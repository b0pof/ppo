package main

import (
	"testing"

	"github.com/b0pof/ppo/tests/db/client"
	"github.com/stretchr/testify/assert"
)

func main() {
	Prepare()
}

func Prepare() {
	db := client.DBConnect()
	defer func() {
		_ = db.Close()
	}()

	// err := client.CleanUpDB(db)
	// assert.NoError(&testing.T{}, err, "db cleanup failed")

	err := client.PrepareDB(db)
	assert.NoError(&testing.T{}, err, "db preparation failed")
}
