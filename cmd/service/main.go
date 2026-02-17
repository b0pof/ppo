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

	"github.com/cohere-ai/cohere-go/v2/client"
	"github.com/cohere-ai/cohere-go/v2/option"
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
	verificationRepo "github.com/b0pof/ppo/internal/repository/verification"
	"github.com/b0pof/ppo/internal/server"
	authUc "github.com/b0pof/ppo/internal/usecase/auth"
	cartUc "github.com/b0pof/ppo/internal/usecase/cart"
	generateDescriptionUc "github.com/b0pof/ppo/internal/usecase/description/generate"
	itemUc "github.com/b0pof/ppo/internal/usecase/item"
	verificationCodeUc "github.com/b0pof/ppo/internal/usecase/notification/post/verification/code"
	orderUc "github.com/b0pof/ppo/internal/usecase/order"
	userUc "github.com/b0pof/ppo/internal/usecase/user"
)

const timeout = 3 * time.Second

//nolint:funlen
func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db := configure.MustInitPostgres(ctx, cfg.Postgres)

	redis := configure.MustInitRedis(cfg.Redis)

	reg := prometheus.NewRegistry()

	// <! Clients
	if len(os.Getenv("LLM_API_KEY")) > 0 {
		fmt.Println("<><><><> APIKEY found in env <><><><>")
	} else {
		fmt.Println(">< >< >< APIKEY found in env >< >< ><")
	}

	llmClientOptions := []option.RequestOption{client.WithToken(os.Getenv("LLM_API_KEY"))}
	llmBaseUrl := os.Getenv("LLM_BASE_URL")
	if len(llmBaseUrl) > 0 {
		llmClientOptions = append(llmClientOptions, client.WithBaseURL(llmBaseUrl))
	}
	llmClient := client.NewClient(llmClientOptions...)
	// !>

	// <! Repositories
	authRepository := authRepo.New(redis, authRepo.WithSessionTTL(cfg.Service.SessionTTL))
	userRepository := userRepo.New(db)
	cartRepository := cartRepo.New(db)
	itemRepository := itemRepo.New(db)
	orderRepository := orderRepo.New(db)
	categoryRepository := categoryRepo.New(db)
	reviewRepository := reviewRepo.New(db)
	verificationRepository := verificationRepo.New(db)
	// !>

	// <! Usecases
	generateDescriptionUsecase := generateDescriptionUc.New(llmClient)
	sendVerificationCodeUsecase := verificationCodeUc.New(cfg.SMTP)
	cartUsecase := cartUc.New(cartRepository)
	authUsecase := authUc.New(authRepository, verificationRepository, userRepository, sendVerificationCodeUsecase)
	itemUsecase := itemUc.New(itemRepository, generateDescriptionUsecase)
	orderUsecase := orderUc.New(orderRepository, itemRepository, cartRepository)
	userUsecase := userUc.New(userRepository, authUsecase)

	authUsecase.SetUserUsecase(userUsecase)
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

	r.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
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
	r.Use(authMiddleware.New(authUsecase, verificationRepository, userUsecase, cfg.Server.Mode))
	r.Use(observabilityMiddleware.New(metrics.NewMetrics(reg), log))
	// r.Use(permsMiddleware.New())
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
	srv := server.NewServer(corsMiddleware.Handler(srvRouter), cfg.Server)
	// !>

	// <! Run
	go func() {
		log.Info(fmt.Sprintf("server is running on port %s...", cfg.Server.Port))
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
