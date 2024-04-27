package app

import (
	"github.com/gin-gonic/gin"
	"github.com/vovk404/course-platform/application-api/config"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/pkg/database"
	"github.com/vovk404/course-platform/application-api/pkg/httpserver"
	"github.com/vovk404/course-platform/application-api/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	sql, err := database.NewPostgreSQL(database.PostgreSQLConfig{
		User:     cfg.PostgreSQL.User,
		Password: cfg.PostgreSQL.Password,
		Host:     cfg.PostgreSQL.Host,
		Database: cfg.PostgreSQL.Database,
	})
	if err != nil {
		log.Fatal("failed to init postgresql", "err", err)
	}

	err = sql.DB.AutoMigrate(
		&entity.User{},
		&entity.Account{},
		&entity.AccountDevices{},
		&entity.AccountSettings{},
	)
	if err != nil {
		log.Fatal("automigration failed", "err", err)
	}

	databases := map[string]database.Database{
		"postgreSQL": sql,
	}

	//services := service2.Services{}

	httpHandler := gin.New()

	//controller.New(&controller.Options{
	//	Handler:  httpHandler,
	//	Services: services,
	//	Logger:   log,
	//	Config:   cfg,
	//})

	httpServer := httpserver.New(
		httpHandler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Second*60),
		httpserver.WriteTimeout(time.Second*60),
		httpserver.ShutdownTimeout(time.Second*30),
	)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify", "err", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown", "err", err)
	}

	for _, db := range databases {
		err = db.Close()
		if err != nil {
			log.Error("app - Run - db.Close", "err", err)
		}
	}
}
