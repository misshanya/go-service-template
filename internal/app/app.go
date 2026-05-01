package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go-service-template/internal/errorz"
	"go-service-template/internal/jwt"
	"go-service-template/internal/transport/http/server"
	v1 "go-service-template/internal/transport/http/v1"
	swaggerui "go-service-template/pkg/swagger-ui"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"github.com/valkey-io/valkey-go"

	"github.com/segmentio/kafka-go"

	"go-service-template/internal/config"
	"go-service-template/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userrepo "go-service-template/internal/repository/user"
	userservice "go-service-template/internal/service/user"
	userhandler "go-service-template/internal/transport/http/v1/user"

	custommiddleware "go-service-template/internal/transport/http/middleware"
)

type App struct {
	cfg          *config.Config
	l            *slog.Logger
	e            *echo.Echo
	dbPool       *pgxpool.Pool
	s3Client     *s3.Client
	kafkaWriter  *kafka.Writer
	kafkaReader  *kafka.Reader
	valkeyClient valkey.Client
}

// New creates and initializes a new instance of App
func New(ctx context.Context, cfg *config.Config, l *slog.Logger) (*App, error) {
	a := &App{
		cfg: cfg,
		l:   l,
	}

	// Your needed connectors below (db, cache, storage, message broker, etc.)

	if err := a.initDB(ctx); err != nil {
		return nil, err
	}

	if err := a.migrateDB(); err != nil {
		return nil, err
	}

	jwtManager := jwt.NewManager(a.cfg.JWT.Secret, time.Hour)

	userRepo := userrepo.New(a.dbPool)
	userService := userservice.New(userRepo, jwtManager)
	userHandler := userhandler.New(userService)

	if err := a.initEcho(); err != nil {
		return nil, err
	}

	apiGroup := a.e.Group("/api/v1")

	strictMiddlewares := []v1.StrictMiddlewareFunc{
		custommiddleware.StrictErrorMiddleware,
	}

	strictServer := server.New(userHandler)
	strictHandler := v1.NewStrictHandler(strictServer, strictMiddlewares)

	oapiValidationMiddleware, err := custommiddleware.OpenAPIValidationMiddleware(jwtManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI validation middleware: %w", err)
	}
	apiGroup.Use(oapiValidationMiddleware)

	v1.RegisterHandlers(apiGroup, strictHandler)

	return a, nil
}

// Start performs a start of all functional services
func (a *App) Start(errChan chan<- error) {
	a.l.Info("Starting...")
	if err := a.e.Start(a.cfg.HttpSrv.Addr); err != nil {
		errChan <- err
	}
}

// Stop performs a graceful shutdown for all components
func (a *App) Stop(ctx context.Context) error {
	a.l.Info("[!] Shutting down...")

	var stopErr error

	a.l.Info("Stopping http server...")
	if err := a.e.Shutdown(ctx); err != nil {
		stopErr = errors.Join(stopErr, fmt.Errorf("failed to shutdown http server: %w", err))
	}

	a.l.Info("Closing database pool...")
	a.dbPool.Close()

	if stopErr != nil {
		return stopErr
	}

	a.l.Info("Stopped gracefully")
	return nil
}

// initDB initializes a new pool for PostgreSQL db
func initDB(ctx context.Context, dbURL string, maxConns int32) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// initDB sets up PostgreSQL db
func (a *App) initDB(ctx context.Context) error {
	dbPool, err := initDB(ctx, a.cfg.Postgres.URL, a.cfg.Postgres.MaxConns)
	if err != nil {
		return fmt.Errorf("failed to init db connection: %w", err)
	}
	a.dbPool = dbPool
	return nil
}

// migrateDB performs a migration to ensure the schema is up to date
func (a *App) migrateDB() error {
	conn := sql.OpenDB(stdlib.GetConnector(*a.dbPool.Config().ConnConfig))
	defer conn.Close()

	return db.Migrate(conn)
}

// initEcho sets up a new Echo instance with logger
func (a *App) initEcho() error {
	a.e = echo.New()
	a.e.HideBanner = true
	a.e.HidePort = true
	a.e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api/v1/swagger")
		},
	}))
	a.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogRemoteIP: true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				a.l.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("ip", v.RemoteIP),
					slog.String("latency", time.Now().Sub(v.StartTime).String()),
				)
			} else {
				a.l.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("ip", v.RemoteIP),
					slog.String("latency", time.Now().Sub(v.StartTime).String()),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	a.e.Use(middleware.Recover())

	a.e.GET("/api/v1/openapi.json", func(c echo.Context) error {
		spec, err := v1.GetSwagger()
		if err != nil {
			slog.Error("failed to get swagger spec", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errorz.ErrInternalServerError.Error())
		}
		return c.JSON(http.StatusOK, spec)
	})

	swaggerUIHandler, err := swaggerui.Handler()
	if err != nil {
		return fmt.Errorf("failed to get swagger ui handler: %w", err)
	}

	uiHandler := http.StripPrefix("/api/v1/swagger", swaggerUIHandler)
	a.e.GET("/api/v1/swagger/*", echo.WrapHandler(uiHandler))
	a.e.GET("/api/v1/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/v1/swagger/")
	})

	return nil
}

// initS3 creates an S3 client and tries to create bucket in the storage
func (a *App) initS3(ctx context.Context) error {
	cfg := aws.Config{
		Region:       a.cfg.S3.Region,
		BaseEndpoint: aws.String(a.cfg.S3.Endpoint),
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				a.cfg.S3.AccessKey,
				a.cfg.S3.SecretKey,
				""),
		),
	}

	a.s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err := a.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(a.cfg.S3.BucketName),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) && apiError.ErrorCode() == "BucketAlreadyOwnedByYou" {
			return nil
		}
		return fmt.Errorf("failed to create s3 bucket: %w", err)
	}

	return nil
}

// initKafkaWriter creates a new Kafka writer with the auto topic creation allowed
func (a *App) initKafkaWriter() {
	a.kafkaWriter = &kafka.Writer{
		Addr:                   kafka.TCP(a.cfg.Kafka.Addr),
		AllowAutoTopicCreation: true,
	}
}

// initKafkaReader creates a new Kafka reader
func (a *App) initKafkaReader() {
	a.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{a.cfg.Kafka.Addr},
		GroupID:     a.cfg.Kafka.ReaderGroupID,
		GroupTopics: []string{a.cfg.Kafka.Topic},
	})
}

// initValkey sets up a connection to cache
func (a *App) initValkey() error {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{a.cfg.Valkey.Addr},
		Password:    a.cfg.Valkey.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to init Valkey connection: %w", err)
	}
	a.valkeyClient = client
	return nil
}
