package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jva44ka/ozon-simulator-go-products/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"

	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/repository"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/service"
	"github.com/jva44ka/ozon-simulator-go-products/internal/infra/config"
	"github.com/jva44ka/ozon-simulator-go-products/internal/infra/http/middlewares"
	"github.com/jva44ka/ozon-simulator-go-products/internal/infra/http/round_trippers"
)

type App struct {
	config *config.Config
	server http.Server
}

type App struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	db         *sql.DB
	cfg        *config.Config
}

func NewApp(configPath string) (*App, error) {
	configImpl, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{
		config: configImpl,
	}

	app.server.Handler, err = boostrapHandler(configImpl)
	if err != nil {
		return nil, fmt.Errorf("boostrapHandler: %w", err)
	}

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return app.server.Serve(l)
}

func boostrapHandler(cfg *config.Config) (http.Handler, error) {
	tr := http.DefaultTransport
	tr = round_trippers.NewTimerRoundTipper(tr)

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	))
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	productRepository := repository.NewPgxRepository(pool)
	productService := service.NewProductService(productRepository)

	mx := http.NewServeMux()
	mx.Handle("GET /product/{sku}", get_product_by_sku_handler.NewGetProductsBySkuHandler(productService))
	mx.Handle("/swagger/", httpSwagger.WrapHandler)

	middleware := middlewares.NewTimerMiddleware(mx)

	return middleware, nil
}
