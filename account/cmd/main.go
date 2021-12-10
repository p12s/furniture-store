package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/p12s/furniture-store/account/internal/broker"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/repository"
	"github.com/p12s/furniture-store/account/internal/service"
	handler "github.com/p12s/furniture-store/account/internal/transport/rest"

	"github.com/sirupsen/logrus"
)

func main() {
	runtime.GOMAXPROCS(1)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error reading env variables from file: %s\n", err.Error())
	}
	cfg, err := config.New()
	if err != nil {
		logrus.Fatalf("error loading env variables: %s\n", err.Error())
	}

	db, err := repository.NewSqlite3DB(repository.Config{Driver: cfg.DB.Driver})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s\n", err.Error())
	}

	repos := repository.NewRepository(db)
	if err != nil {
		logrus.Fatalf("failed to initialize authentication token ttl: %s\n", err.Error())
	}
	services := service.NewService(repos, &cfg.Auth)
	broker, err := broker.NewBroker(services, &cfg.Broker)
	if err != nil {
		logrus.Fatalf("kafka error: %s\n", err.Error())
	}
	handlers := handler.NewHandler(services, broker)

	srv := new(Server)
	go func() {
		if err := srv.Run(cfg.Server.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error while running http server: %s\n", err.Error())
		}
	}()
	logrus.Print("😀 account app started with port: ", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("account app shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}
	// TODO broker close
}

// Server - http server
type Server struct {
	httpServer *http.Server
}

// Run - start
func (s *Server) Run(port int, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown - grace-full
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
