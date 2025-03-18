package main

import (
	"net/http"

	"github.com/Axel791/order/internal/config"
	"github.com/Axel791/order/internal/db"
	"github.com/Axel791/order/internal/grpc/v1/pb"
	apiV1Handlers "github.com/Axel791/order/internal/rest/v1"
	"github.com/Axel791/order/internal/usecases/order/repositories"
	"github.com/Axel791/order/internal/usecases/order/scenarios"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(logrus.InfoLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	log.Infof("databse_dsn: %s", cfg.DatabaseDSN)
	dbConn, err := db.ConnectDB(cfg.DatabaseDSN, cfg)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer func() {
		if dbConn != nil {
			_ = dbConn.Close()
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)

	// gRPC init
	conn, err := grpc.NewClient(
		cfg.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create new client: %v", err)
	}

	defer conn.Close()

	// repositories
	orderRepository := repositories.NewSqlOrderRepository(dbConn)

	// gRPC Clients
	loyaltyClient := pb.NewConclusionUserBalanceUseCaseClient(conn)

	// use cases
	createOrderUseCase := scenarios.NewCreateOrderUseCase(orderRepository, loyaltyClient)

	// rest
	router.Route("/api/v1", func(r chi.Router) {
		r.Method(
			http.MethodPost,
			"order/",
			apiV1Handlers.NewCreateOrderHandler(log, createOrderUseCase),
		)
	})

	log.Infof("server started on %s", cfg.Address)
	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

}
