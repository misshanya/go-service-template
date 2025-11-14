package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go-service-template/internal/config"
	"go-service-template/internal/db"
	"go-service-template/internal/db/sqlc/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userrepo "go-service-template/internal/repository/user"
	userservice "go-service-template/internal/service/user"
	userhandler "go-service-template/internal/transport/http/v1/user"
)

type App struct {
	cfg    *config.Config
	l      *slog.Logger
	e      *echo.Echo
	dbPool *pgxpool.Pool
}

// New creates and initializes a new instance of App
func New(ctx context.Context, cfg *config.Config, l *slog.Logger) (*App, error) {
	a := &App{
		cfg: cfg,
		l:   l,
	}

	if err := a.initDB(ctx); err != nil {
		return nil, err
	}

	if err := a.migrateDB(); err != nil {
		return nil, err
	}

	queries := storage.New(a.dbPool)
	userRepo := userrepo.New(queries)
	userService := userservice.New(userRepo)
	userHandler := userhandler.New(userService)

	a.initEcho()

	apiGroup := a.e.Group("/api/v1")
	userGroup := apiGroup.Group("/user")

	userHandler.Setup(userGroup)

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

	a.l.Info("Closing database pool...")
	a.dbPool.Close()

	a.l.Info("Stopping http server...")
	if err := a.e.Shutdown(ctx); err != nil {
		stopErr = errors.Join(stopErr, fmt.Errorf("failed to shutdown http server: %w", err))
	}

	if stopErr != nil {
		return stopErr
	}

	a.l.Info("Stopped gracefully")
	return nil
}

// initDB initializes a new pool for PostgreSQL db
func initDB(ctx context.Context, dbURL string, maxConns int32) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = maxConns

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
	return db.Migrate(sql.OpenDB(stdlib.GetConnector(*a.dbPool.Config().ConnConfig)))
}

// initEcho sets up a new Echo instance with logger
func (a *App) initEcho() {
	a.e = echo.New()
	a.e.HideBanner = true
	a.e.HidePort = true
	a.e.Pre(middleware.RemoveTrailingSlash())
	a.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
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
}
