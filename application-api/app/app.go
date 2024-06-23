package app

import (
	"github.com/gin-gonic/gin"
	"github.com/vovk404/course-platform/application-api/config"
	controller "github.com/vovk404/course-platform/application-api/internal/controller/http"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/internal/service"
	"github.com/vovk404/course-platform/application-api/internal/storage"
	"github.com/vovk404/course-platform/application-api/pkg/auth"
	"github.com/vovk404/course-platform/application-api/pkg/database"
	"github.com/vovk404/course-platform/application-api/pkg/hash"
	"github.com/vovk404/course-platform/application-api/pkg/httpserver"
	"github.com/vovk404/course-platform/application-api/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	// Wait until Postgres started
	time.Sleep(2 * time.Second)
	sql := connectToDB(cfg, log)

	err := sql.DB.AutoMigrate(
		&entity.User{},
		&entity.Account{},
		&entity.AccountDevices{},
		&entity.AccountSettings{},
	)
	if err != nil {
		log.Fatal("automigration failed", "err", err)
	}

	storages := service.Storages{
		UserStorage:    storage.NewUserStorage(sql),
		AccountStorage: storage.NewAccountStorage(sql),
		NodeStorage:    storage.NewNodeStorage(sql),
	}

	databases := map[string]database.Database{
		"postgreSQL": sql,
	}

	serviceOptions := &service.Options{
		Storages: &storages,
		Config:   cfg,
		Logger:   log,
		Hash:     hash.NewHash(),
		Auth:     auth.NewAuth(),
	}

	services := service.Services{
		AuthService:    service.NewAuthService(serviceOptions),
		AccountService: service.NewAccountService(serviceOptions),
		NodeService:    service.NewNodeService(serviceOptions),
	}

	httpHandler := gin.New()

	controller.New(&controller.Options{
		Handler:  httpHandler,
		Services: services,
		Logger:   log,
		Config:   cfg,
	})

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

	// Shut down server after 30 sec (according to httpserver.ShutdownTimeout(time.Second*30))
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

// Try to connect to Postgress 10 times before throwing an error, postgress starting after the application-api.
func connectToDB(cfg *config.Config, logger logger.Logger) *database.PostgreSQL {
	var counts int

	for {
		connection, err := database.NewPostgreSQL(database.PostgreSQLConfig{
			User:     cfg.PostgreSQL.User,
			Password: cfg.PostgreSQL.Password,
			Host:     cfg.PostgreSQL.Host,
			Database: cfg.PostgreSQL.Database,
			Port:     cfg.PostgreSQL.Port,
		})
		if err != nil {
			logger.Error("Postgres is not ready")
			counts++
		} else {
			logger.Info("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			logger.Fatal(err.Error())
		}
		logger.Debug("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
