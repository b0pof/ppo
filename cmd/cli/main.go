package main

import (
	"context"
	"log"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/cli"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/config"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/configure"
	authRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/auth"
	cartRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/cart"
	itemRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/item"
	orderRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/order"
	userRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/user"
	authUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/auth"
	cartUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/cart"
	itemUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/item"
	orderUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/order"
	userUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/user"
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
