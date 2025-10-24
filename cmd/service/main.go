package main

import (
	"context"
	"fmt"
	authMiddleware "git.iu7.bmstu.ru/kia22u475/ppo/internal/middleware/auth"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	apiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler"
	authApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/auth"
	cartApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/cart"
	categoryApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/category"
	itemApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/item"
	orderApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/order"
	reviewApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/review"
	sellerApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/seller"
	userApiHandler "git.iu7.bmstu.ru/kia22u475/ppo/internal/api/handler/user"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/config"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/configure"
	sdk "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	permissionMiddleware "git.iu7.bmstu.ru/kia22u475/ppo/internal/middleware/permission"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	authRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/auth"
	cartRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/cart"
	categoryRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/category"
	itemRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/item"
	orderRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/order"
	reviewRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/review"
	userRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/user"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/server"
	authUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/auth"
	cartUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/cart"
	itemUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/item"
	orderUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/order"
	userUc "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/user"
)

const timeout = 3 * time.Second

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db := configure.MustInitPostgres(ctx, cfg.Postgres)

	redis := configure.MustInitRedis(cfg.Redis)

	// <! Repositories
	authRepository := authRepo.New(redis, authRepo.WithSessionTTL(cfg.Service.SessionTTL))
	userRepository := userRepo.New(db)
	cartRepository := cartRepo.New(db)
	itemRepository := itemRepo.New(db)
	orderRepository := orderRepo.New(db)
	categoryRepository := categoryRepo.New(db)
	reviewRepository := reviewRepo.New(db)
	// !>

	// <! Usecases
	authUsecase := authUc.New(authRepository, userRepository)
	cartUsecase := cartUc.New(cartRepository)
	itemUsecase := itemUc.New(itemRepository)
	orderUsecase := orderUc.New(orderRepository, itemRepository, cartRepository)
	userUsecase := userUc.New(userRepository)
	// !>

	// <! Router
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `Not found`, 404)
	})
	// !>

	// <! Permissions
	permsMiddleware := permissionMiddleware.New()

	for path, perms := range model.Resources {
		permsMiddleware.Register(path, perms)
	}
	// !>

	// <! Middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	r.Use(authMiddleware.New(authUsecase, userUsecase))
	r.Use(permsMiddleware.New())
	// !>

	// <! Handlers
	authHandler := authApiHandler.New(authUsecase, log, cfg.Service.SessionTTL)
	cartHandler := cartApiHandler.New(cartUsecase, log)
	categoryHandler := categoryApiHandler.New(categoryRepository, itemRepository, log)
	itemHandler := itemApiHandler.New(itemUsecase, log)
	orderHandler := orderApiHandler.New(orderUsecase, log)
	reviewHandler := reviewApiHandler.New(reviewRepository, log)
	sellerHandler := sellerApiHandler.New(itemUsecase, log)
	userHandler := userApiHandler.New(userUsecase, log)
	handlerPerformer := apiHandler.NewHandler(
		authHandler,
		cartHandler,
		categoryHandler,
		itemHandler,
		orderHandler,
		reviewHandler,
		sellerHandler,
		userHandler,
	)
	srvRouter := sdk.HandlerFromMux(handlerPerformer, r)
	// !>

	// <! Server
	srv := server.NewServer(corsMiddleware.Handler(srvRouter))
	// !>

	// <! Run
	go func() {
		log.Info("server is running...")
		if err := srv.Run(); err != nil {
			log.Error(">>> ERROR: HTTP server ListenAndServe error: " + err.Error())
		}
	}()
	// !>

	// <! Graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	log.Info("shutting down...")
	if err := srv.Stop(ctx); err != nil {
		log.Error(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	// !>
}
