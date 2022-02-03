package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daniilty/secrets-manager/internal/core"
	"github.com/daniilty/secrets-manager/internal/db"
	"github.com/daniilty/secrets-manager/internal/healthcheck"
	"github.com/daniilty/secrets-manager/internal/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	// set more descriptive exit codes...
	exitCodeNotOK = 2
)

func run() (int, error) {
	cfg, err := loadEnvConfig()
	if err != nil {
		return exitCodeNotOK, err
	}

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		return exitCodeNotOK, err
	}

	sugared := logger.Sugar()

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.mongoConnString))
	if err != nil {
		return exitCodeNotOK, err
	}

	mongoDB := mongoClient.Database(cfg.mongoDBName)
	mongoCollection := mongoDB.Collection(cfg.mongoCollectionName)

	pinger := db.NewMongoPinger(db.WithMongoClient(mongoClient))
	appInfo := healthcheck.NewChecker(healthcheck.WithMongoDBPinger(pinger))

	dbImpl := db.NewMongoDB(db.WithMongoCollection(mongoCollection))
	service := core.NewSecretsManager(core.WithDB(dbImpl))

	http.DefaultServeMux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		bb, err := json.Marshal(appInfo.Check())
		if err != nil {
			sugared.Errorw("Health check.", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

			return
		}

		w.Write(bb)
	})
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	devopsServer := &http.Server{
		Addr:    cfg.httpDevopsAddr,
		Handler: http.DefaultServeMux,
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		sugared.Infow("Server started.", "server", "devops", "addr", devopsServer.Addr)
		err := devopsServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			sugared.Errorw("Listen and serve.", "server", "devops", "err", err)
		}
	}()

	secretsServer := server.NewHTTP(cfg.httpSecretsAddr, sugared, service)

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		secretsServer.Run(ctx)
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINFO, syscall.SIGTERM, os.Interrupt)

	<-term
	cancel()
	sugared.Infow("Server shutdown.", "server", "devops")
	devopsServer.Shutdown(context.Background())
	wg.Wait()

	return 0, nil
}

func main() {
	code, err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(code)
	}
}
