// Package main creates and runs application instance.
package main

import (
	"github.com/vovk404/course-platform/application-api/app"
	"github.com/vovk404/course-platform/application-api/config"
	"github.com/vovk404/course-platform/application-api/pkg/logger"
)

func main() {

	log := logger.New("main")

	cfg := config.Get()
	log.Info("read config", "config", cfg)

	app.Run(cfg)
}
