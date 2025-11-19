package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	apiHandler "github.com/b0pof/ppo/internal/api/handler"
	authApiHandler "github.com/b0pof/ppo/internal/api/handler/auth"
	cartApiHandler "github.com/b0pof/ppo/internal/api/handler/cart"
	categoryApiHandler "github.com/b0pof/ppo/internal/api/handler/category"
	itemApiHandler "github.com/b0pof/ppo/internal/api/handler/item"
	orderApiHandler "github.com/b0pof/ppo/internal/api/handler/order"
	reviewApiHandler "github.com/b0pof/ppo/internal/api/handler/review"
	sellerApiHandler "github.com/b0pof/ppo/internal/api/handler/seller"
	userApiHandler "github.com/b0pof/ppo/internal/api/handler/user"
	"github.com/b0pof/ppo/internal/config"
	"github.com/b0pof/ppo/internal/configure"
	sdk "github.com/b0pof/ppo/internal/generated"
	authMiddleware "github.com/b0pof/ppo/internal/middleware/auth"
	observabilityMiddleware "github.com/b0pof/ppo/internal/middleware/observability"
	permissionMiddleware "github.com/b0pof/ppo/internal/middleware/permission"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/pkg/metrics"
	authRepo "github.com/b0pof/ppo/internal/repository/auth"
	cartRepo "github.com/b0pof/ppo/internal/repository/cart"
	categoryRepo "github.com/b0pof/ppo/internal/repository/category"
	itemRepo "github.com/b0pof/ppo/internal/repository/item"
	orderRepo "github.com/b0pof/ppo/internal/repository/order"
	reviewRepo "github.com/b0pof/ppo/internal/repository/review"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	"github.com/b0pof/ppo/internal/server"
	authUc "github.com/b0pof/ppo/internal/usecase/auth"
	cartUc "github.com/b0pof/ppo/internal/usecase/cart"
	itemUc "github.com/b0pof/ppo/internal/usecase/item"
	orderUc "github.com/b0pof/ppo/internal/usecase/order"
	userUc "github.com/b0pof/ppo/internal/usecase/user"
)

const timeout = 3 * time.Second

////go:embed api/schema.yml
//var spec embed.FS

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db := configure.MustInitPostgres(ctx, cfg.Postgres)

	redis := configure.MustInitRedis(cfg.Redis)

	reg := prometheus.NewRegistry()

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
	// !!>

	specData, err := os.ReadFile("api/schema.yml")
	if err != nil {
		panic(err)
	}

	r.Handle("/api/schema.yml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(specData)
	}))

	opts := apiMiddleware.SwaggerUIOpts{
		Path:    "/docs",
		SpecURL: "/api/schema.yml",
		Title:   "API Documentation",
	}

	r.Handle("/docs/", apiMiddleware.SwaggerUI(opts, nil))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		http.Error(w, `Not found!`, 404)
	})

	r.Handle("/public/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			Registry: reg,
		},
	))
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
	r.Use(observabilityMiddleware.New(metrics.NewMetrics(reg), log))
	//r.Use(permsMiddleware.New())
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
