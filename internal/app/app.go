package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jva44ka/ozon-simulator-go-products/internal/app/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "github.com/jva44ka/ozon-simulator-go-products/internal/app/gen/ozon-simulator-go-products/api/v1/proto"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/repository"
	"github.com/jva44ka/ozon-simulator-go-products/internal/domain/service"
	"github.com/jva44ka/ozon-simulator-go-products/internal/infra/config"
)

type repositoryMetrics struct {
	requestsTotal          *prometheus.CounterVec
	optimisticLockFailures prometheus.Counter
}

func newRepositoryMetrics() *repositoryMetrics {
	return &repositoryMetrics{
		requestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "products_db_requests_total",
			Help: "Total number of products DB requests",
		}, []string{"method", "status"}),
		optimisticLockFailures: promauto.NewCounter(prometheus.CounterOpts{
			Name: "products_optimistic_lock_failures_total",
			Help: "Total number of optimistic lock failures in product count updates",
		}),
	}
}

func (m *repositoryMetrics) RecordRequest(method, status string) {
	m.requestsTotal.WithLabelValues(method, status).Inc()
}

func (m *repositoryMetrics) IncOptimisticLockFailure() {
	m.optimisticLockFailures.Inc()
}

type App struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	cfg        *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
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

	repo := repository.NewProductRepository(pool, newRepositoryMetrics())
	domainService := service.NewProductService(repo)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Panic,
			middleware.Metrics,
			middleware.Logger,
			middleware.Auth(cfg),
			middleware.Validate,
		),
	)
	grpcService := NewGrpcService(domainService)

	pb.RegisterProductsServer(grpcServer, grpcService)

	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(
			func(header string) (string, bool) {
				switch strings.ToLower(header) {
				case "x-auth":
					return header, true
				default:
					return runtime.DefaultHeaderMatcher(header)
				}
			},
		))

	err = pb.RegisterProductsHandlerFromEndpoint(
		ctx,
		mux,
		cfg.GrpcServer.Host+":"+cfg.GrpcServer.Port,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		return nil, err
	}

	reflection.Register(grpcServer)

	httpMux := http.NewServeMux()
	// grpc gateway
	httpMux.Handle("/", mux)
	// swagger json
	httpMux.Handle("/api/", http.StripPrefix(
		"/api/",
		http.FileServer(http.Dir("./swagger/api/v1")),
	))
	// swagger UI
	httpMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprint(w, swaggerUiHtml)
	})
	// prometheus metrics
	httpMux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    cfg.HttpServer.Host + ":" + cfg.HttpServer.Port,
		Handler: httpMux,
	}

	return &App{
		grpcServer: grpcServer,
		httpServer: httpServer,
		cfg:        cfg,
	}, nil
}

func (a *App) Run() error {

	lis, err := net.Listen("tcp", ":"+a.cfg.GrpcServer.Port)
	if err != nil {
		return err
	}

	go func() {
		a.grpcServer.Serve(lis)
	}()

	return a.httpServer.ListenAndServe()
}
