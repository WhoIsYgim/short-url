package app

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"short-link/config"
	grpcHandlers "short-link/internal/delivery/grpc/handlers"
	"short-link/internal/delivery/http/handlers"
	"short-link/internal/delivery/http/middleware"
	"short-link/internal/repository/memory"
	pgRepo "short-link/internal/repository/postgres"
	"short-link/internal/repository/tokens"
	"short-link/internal/usecase"
	"short-link/pkg/connector/postgres"
	"short-link/pkg/grpc/api"
	"short-link/pkg/link/generator"
	"strconv"
	"syscall"
)

var DbType string

func init() {
	flag.StringVar(&DbType, "db", "pg", `Type of database for storing short links.	
		Supported options: 
			"pg"     for PostrgreSQL
			"im"  	 for in-memory solution
	`)
}

func Run() {
	flag.Parse()
	summary, err := config.ParseConfig()
	if err != nil {
		logrus.Fatal("error with parsing config")
	}
	cfg := &summary.Cfg
	dbCfg := &summary.DbCfg
	var storage usecase.LinkRepository

	switch DbType {
	case "pg":
		pgDefault, err := postgres.GetPostgresConnector(dbCfg)
		if err != nil {
			logrus.Fatal(err)
		}
		sqlxDB := postgres.GetSqlxConnector(pgDefault, "postgres")
		storage = pgRepo.NewLinkStorage(sqlxDB)
	case "im":
		storage = memory.NewLinkStorage()
	default:
		flag.Usage()
		logrus.Fatal("incorrect db option", DbType)
	}

	tokenCache := tokens.NewTokenCache()
	gd := generator.GenData{
		Alphabet: cfg.LinkConfig.Alphabet,
		Length:   cfg.LinkConfig.TokenLength,
	}
	strGenerator := generator.NewGeneratorWithData(gd)

	service := usecase.NewLinkService(storage, cfg, tokenCache, strGenerator)

	handler := handlers.NewShortLinkHandler(service)
	router := InitRouter(handler)

	server := http.Server{
		Handler: router,
		Addr:    ":" + strconv.FormatUint(uint64(cfg.ServiceConfig.Port), 10),
	}

	go func() {
		grpcHandler := grpcHandlers.NewLinkHandler(service)

		srv := grpc.NewServer()
		api.RegisterShortLinkServiceServer(srv, grpcHandler)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServiceConfig.GrpcPort))
		if err != nil {
			logrus.Fatal(err)
		}
		if err := srv.Serve(listener); err != nil {
			logrus.Fatal(err)
		}
	}()

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			logrus.Fatal("start server error:", err)
			return
		}
	}()

	logrus.Info("service statred")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	err = storage.ShutDown()
	if err != nil {
		logrus.Fatal("DB shutdown with error: ", err)
	}
	logrus.Info("shutdown")
}

func InitRouter(handler *handlers.ShortLinkHandler) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ErrorMiddleware())

	router.GET("/url/:key", handler.GetLink)
	router.POST("/url/", handler.CreateLink)

	return router
}
