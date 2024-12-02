package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aeilang/urlshortener/config"
	"github.com/aeilang/urlshortener/db"
	apiv1 "github.com/aeilang/urlshortener/internal/api/v1"
	"github.com/aeilang/urlshortener/internal/cache"
	"github.com/aeilang/urlshortener/internal/service"
	"github.com/aeilang/urlshortener/pkg/shortcode"
	"github.com/aeilang/urlshortener/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	e          *echo.Echo
	cfg        *config.Config
	db         *sql.DB
	cache      cache.Cache
	generator  shortcode.ShortCodeGenerator
	urlService service.URLService
	urlHandler *apiv1.URLHandler
}

func NewApplication() (*Application, error) {
	a := &Application{}
	if err := a.loadConfig(); err != nil {
		return nil, fmt.Errorf("loadConfig: %w", err)
	}
	if err := a.initDB(); err != nil {
		return nil, fmt.Errorf("initdb: %w", err)
	}
	if err := a.initCache(); err != nil {
		return nil, fmt.Errorf("initCache: %w", err)
	}
	a.initGenerator()
	a.initHandler()
	a.initEcho()
	a.addRoutes()
	return a, nil
}

func (a *Application) Run() {
	go a.clearnup()
	go a.startServer()
	a.shudown()
}

func (a *Application) loadConfig() error {
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		return err
	}

	a.cfg = cfg
	return nil
}

func (a *Application) initDB() error {
	sqldb, err := db.InitDB(a.cfg.Database)
	if err != nil {
		return err
	}

	a.db = sqldb
	return nil
}

func (a *Application) initCache() error {
	cache, err := cache.NewRedisCache(a.cfg.Redis)
	if err != nil {
		return err
	}
	a.cache = cache
	return nil
}

func (a *Application) initGenerator() {
	a.generator = shortcode.NewShortCodeGenerator(a.cfg.ShortCode.MinLength)
}

func (a *Application) initHandler() {
	a.urlService = service.NewURLService(a.db, a.cache, a.generator, a.cfg)
	a.urlHandler = apiv1.NewURLHandler(a.urlService, a.cfg.App.BaseURL)
}

func (a *Application) initEcho() {
	e := echo.New()
	e.Validator = validator.NewCustomeValidator()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	a.e = e
}

func (a *Application) addRoutes() {
	a.e.POST("/api/url", a.urlHandler.CreateURL)
	a.e.GET("/:code", a.urlHandler.RedirectURL)
}

func (a *Application) clearnup() {
	ticker := time.NewTicker(a.cfg.App.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := a.urlService.Cleanup(context.Background()); err != nil {
			log.Printf("Failed to clean expired URLs: %v", err)
		}
	}
}

func (a *Application) startServer() {
	if err := a.e.Start(a.cfg.Server.Address); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}

func (a *Application) shudown() {
	// Wait for interupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer a.db.Close()
	defer a.cache.Close()
	// Gracefull shudown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
