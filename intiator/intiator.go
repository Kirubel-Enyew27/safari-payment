package intiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/handler/middleware"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/persistencedb"
	"github.com/Kirubel-Enyew27/safari-payment/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Intiate() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf(`{"level":"fatal","msg":"failed to initialize logger: %v"}
`, err)
		os.Exit(1)
	}

	logger.Info("initializing config")
	err = godotenv.Load("./config/.env")
	if err != nil {
		logger.Fatal("unable to intialize config", zap.Error(err))
	}
	logger.Info("config initialized")

	logger.Info("initializing database")
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		logger.Fatal("unable to get DATABASE_URL")
	}
	pgxConn := InitDB(db_url, logger)
	logger.Info("database initialized")

	if strings.ToLower(os.Getenv("MIGRATION_ACTIVE")) == "true" {
		logger.Info("initializing migration")
		migration_path := os.Getenv("MIGRATION_PATH")
		if migration_path == "" {
			logger.Fatal("unable to get MIGRATION_PATH")
		}
		m := InitiateMigration(migration_path, db_url, logger)
		UpMigration(m, logger)
		logger.Info("migration initialized")
	}

	logger.Info("initializing storage layer")
	storage := InitPersistence(persistencedb.New(pgxConn, logger), logger)
	logger.Info("storage layer initialized")

	logger.Info("initializing service layer")
	token, err := utils.GetSafariAccessToken()
	service := InitService(storage, token, logger)
	logger.Info("service layer initialized")

	logger.Info("initializing handler")
	serverTimeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		logger.Fatal("unable to parse server timeout duration", zap.Error(err))
	}
	handler := InitHandler(service, logger, serverTimeout)
	logger.Info("handler initialized")

	logger.Info("initializing server")
	server := gin.New()
	gin.SetMode(gin.DebugMode)
	server.Use(middleware.CORSMiddleware(), middleware.ErrorHandler())
	logger.Info("server initialized")

	logger.Info("initializing router")
	v1 := server.Group("/v1")
	InitRouter(v1, handler, service, logger)
	logger.Info("router initialized")

	readHeaderTimeout, err := time.ParseDuration(os.Getenv("SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		logger.Fatal("unable to parse read header timeout duration", zap.Error(err))
	}

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	if host == "" || port == "" {
		logger.Fatal("unable to get host or port", zap.String("host", host), zap.String("port", port))
	}
	srv := &http.Server{
		Addr:              os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           server,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverPort, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatal("unable to parse server port", zap.Error(err))
	}
	logger.Info("server started",
		zap.String("host", host),
		zap.Int("port", serverPort),
		zap.Time("start_time", time.Now()))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("error listening to the server", zap.Error(err))
		}
	}()

	// Wait for termination signal
	sig := <-quit
	logger.Info("server shutting down", zap.String("signal", sig.String()))

	if serverTimeout == 0 {
		serverTimeout = 5 * time.Second // Default to 5 seconds
		logger.Warn("server timeout not set, using default", zap.Duration("timeout", serverTimeout))
	}
	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("error while shutting down server", zap.Error(err))
	}

	logger.Info("server shutdown complete")

}
