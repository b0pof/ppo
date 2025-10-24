package client

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/config"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/configure"
)

const (
	maxAttempts = 3
	delay       = time.Millisecond * 100
	connTimeout = 10 * time.Second
)

type state struct {
	fixtures *testfixtures.Loader
	tables   []string
}

var s = &state{}

func DbConnect() *sqlx.DB {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	return configure.MustInitPostgres(ctx, config.PostgresConfig{
		Host:     "localhost", // "postgres-ppo-e2e",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		Sslmode:  "disable",
	})
}

func RedisConnect() *redis.Client {
	return configure.MustInitRedis(config.RedisConfig{
		Addr: "localhost:6380",
	})
}

func cleanTable(table string, tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("ALTER TABLE \"%s\" DISABLE TRIGGER ALL;", table))
	if err != nil {
		return fmt.Errorf("failed to disable trigger: %w", err)
	}
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM \"%s\" WHERE TRUE;", table))
	if err != nil {
		return fmt.Errorf("failed to exec delete: %w", err)
	}
	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE \"%s\" ENABLE TRIGGER ALL;", table))
	if err != nil {
		return fmt.Errorf("failed to enable trigger: %w", err)
	}
	return nil
}

func cleanUpDB(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	if len(s.tables) == 0 {
		var tables []string
		err := tx.Select(&tables, "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';")
		if err != nil {
			return err
		}
		s.tables = tables
	}

	for _, table := range s.tables {
		err = cleanTable(table, tx)
		if err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func CleanUpDB(t provider.T, db *sqlx.DB) error {
	var err error
	name := t.Name()
	start := time.Now()
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = cleanUpDB(db)
		if err != nil {
			t.Logf("Failed to cleanup for test %s, attempt %d: %s", name, attempt, err)
			time.Sleep(delay * time.Duration(attempt))
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed to exec cleanup: %w", err)
	}
	t.Logf("Cleanup for test %s finished in %s", name, time.Since(start))

	return nil
}

func CleanUpRedis(t provider.T, r *redis.Client) {
	r.FlushAll()
}

func prepareDB(db *sqlx.DB) error {
	sqlDB := db.DB

	if s.fixtures == nil {
		fixtures, err := testfixtures.New(
			testfixtures.DangerousSkipTestDatabaseCheck(),
			testfixtures.Database(sqlDB),
			testfixtures.Dialect("postgres"),
			testfixtures.Directory("../../db/fixtures"),
			testfixtures.ResetSequencesTo(10000+rand.Int63n(10000)),
		)
		if err != nil {
			return err
		}
		s.fixtures = fixtures
	}

	return s.fixtures.Load()
}

func PrepareDB(t provider.T, db *sqlx.DB) error {
	var err error
	name := t.Name()
	start := time.Now()
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = prepareDB(db)
		if err != nil {
			t.Logf("Failed to prepare DB for test %s, attempt %d: %s", name, attempt, err)
			time.Sleep(delay * time.Duration(attempt))
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed to prepare db: %w", err)
	}
	t.Logf("DB preparation for test %s finished in %s", name, time.Since(start))

	return nil
}
