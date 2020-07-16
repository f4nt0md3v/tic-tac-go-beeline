package app

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/constants"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/db/postgres"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/handlers"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/env"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/errorx"
)

func StartApplication() {
	e := echo.New()

	// Specify middleware to use
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("./public"))

	// Specify handlers
	e.GET("/ws", handlers.WebsocketHandler)

	// Setup Postgres
	port, err := strconv.Atoi(env.GetEnvVariable(constants.DbPort))
	errorx.Must(err)
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

	// Migrate database if not exists
	err = postgres.Migrate(db)
	errorx.Must(err)

	// os.Exit(1) in case we can't start from the port
	e.Logger.Fatal(e.Start(":1323"))
}
