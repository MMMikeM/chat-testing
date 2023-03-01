//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -generate "types,server" -o openapi/openapi.gen.go --package openapi ../docs/swagger.yaml
package internal

import (
	"context"
	"messanger/internal/openapi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Server *echo.Echo
	DB     *gorm.DB
}

func NewApp(ctx context.Context) *App {
	a := App{}

	logger := logrus.New()

	dsn := "host=postgres user=postgres password=postgres dbname=messager port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db

	server := NewServer(ctx, logger, db)

	// Migrate the schema
	// This wont be done this way in production
	db.AutoMigrate(Conversation{})
	db.AutoMigrate(User{})
	db.AutoMigrate(Message{})

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/docs", "docs")
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/ws", server.ws)

	openapi.RegisterHandlers(e, server)
	a.Server = e

	return &a
}

func (a *App) Run() error {
	return a.Server.Start("0.0.0.0:3000")
}
