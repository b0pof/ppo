package main

import (
	"context"
	"log"

	"github.com/b0pof/ppo/internal/cli/cli"
	"github.com/b0pof/ppo/internal/config"
	"github.com/b0pof/ppo/internal/configure"
	authRepo "github.com/b0pof/ppo/internal/repository/auth"
	cartRepo "github.com/b0pof/ppo/internal/repository/cart"
	itemRepo "github.com/b0pof/ppo/internal/repository/item"
	orderRepo "github.com/b0pof/ppo/internal/repository/order"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	authUc "github.com/b0pof/ppo/internal/usecase/auth"
	cartUc "github.com/b0pof/ppo/internal/usecase/cart"
	itemUc "github.com/b0pof/ppo/internal/usecase/item"
	orderUc "github.com/b0pof/ppo/internal/usecase/order"
	userUc "github.com/b0pof/ppo/internal/usecase/user"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := configure.MustInitPostgres(ctx, cfg.Postgres)

	redis := configure.MustInitRedis(cfg.Redis)

	// <! Repositories
	authRepository := authRepo.New(redis, authRepo.WithSessionTTL(cfg.Service.SessionTTL))
	userRepository := userRepo.New(db)
	cartRepository := cartRepo.New(db)
	itemRepository := itemRepo.New(db)
	orderRepository := orderRepo.New(db)
	// !>

	// <! Usecases
	authUsecase := authUc.New(authRepository, userRepository)
	cartUsecase := cartUc.New(cartRepository)
	itemUsecase := itemUc.New(itemRepository)
	orderUsecase := orderUc.New(orderRepository, itemRepository, cartRepository)
	userUsecase := userUc.New(userRepository)
	// !>

	// CLI
	cliApp := cli.New(userUsecase, authUsecase, cartUsecase, orderUsecase, itemUsecase)

	// run
	if err := cliApp.Execute(); err != nil {
		log.Fatal(err)
	}

	// graceful
}
