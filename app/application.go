package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/constants"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/db/postgres"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/handlers"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/context"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/env"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/errorx"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/middleware"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/middleware/logx"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/websocketx"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/repositories"
)

func StartApplication() {
	// Initialize logger
	logger, _ := zap.NewDevelopment()
	defer func() {
		if err := logger.Sync(); err != nil { // flushes buffer, if any
			log.Println(err)
		}
	}()
	sugar := logger.Sugar()

	// Setup Postgres
	port, err := strconv.Atoi(env.GetEnvVariable(constants.DbPort))
	errorx.Must(err)
	sugar.Infof("Connecting to Postgres database")
	db, err := postgres.NewDBSession(postgres.Config{
		Host:             env.GetEnvVariable(constants.DbHost),
		Port:             port,
		User:             env.GetEnvVariable(constants.DbUser),
		Password:         env.GetEnvVariable(constants.DbPass),
		Database:         env.GetEnvVariable(constants.DbName),
		Params:           "",
		ConnectionString: "",
		Mode:             env.GetAppMode(),
	})
	errorx.Must(err)

	sugar.Infof("Migrating data to the database")
	// Migrate database if not exists
	err = postgres.Migrate(db)
	errorx.Must(err)

	sugar.Infof("Database is ready to go")
	repo := repositories.NewGameRepo(db)
	dbCtx := context.NewDbContext(db)

	sugar.Infof("WebSocket pool connections initialized")
	// Creating new pool of ws connections for concurrent writes
	pool := websocketx.NewPool()
	go pool.Run()

	sugar.Infof("Application Context initialized")
	// Setup application context to inject in later into handlers
	appCtx := context.NewAppContext(dbCtx, repo, sugar, pool)

	sugar.Infof("Setting up endpoints and handlers")
	// Setup file server for static frontend files
	http.Handle("/", http.FileServer(http.Dir("./public/")))
	// Setup websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(appCtx, w, r)
	})
	// Specify health check endpoint
	http.HandleFunc("/health", middleware.Chain(handlers.HealthHandler, logx.Logger()))

	// Fire up the server on port APP_PORT
	appPort := fmt.Sprintf(":%s", env.GetEnvVariable(constants.AppPort))
	sugar.Infof("Starting application at port: %s", appPort)
	err = http.ListenAndServe(appPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
