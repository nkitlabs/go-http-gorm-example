package main

import (
	"fmt"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	bookservice "github.com/nkitlabs/go-http-gorm-example/pkg/books/service"
	booktypes "github.com/nkitlabs/go-http-gorm-example/pkg/books/types"
	"github.com/nkitlabs/go-http-gorm-example/pkg/config"
	dbstore "github.com/nkitlabs/go-http-gorm-example/pkg/db"
	"github.com/nkitlabs/go-http-gorm-example/pkg/middleware"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	logger := zap.Must(zap.NewProduction())
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()

	config_filename := os.Getenv("CONFIG_FILE")
	if config_filename == "" {
		config_filename = "config.dev.yaml"
	}

	logger.Info(fmt.Sprintf("Read config file from %s", config_filename))
	conf, err := config.ReadConfig(config_filename)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	db, err := dbstore.Init(conf.Conn)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if err := db.AutoMigrate(&booktypes.Book{}); err != nil {
		logger.Error(err.Error())
		return
	}

	bookRepository := bookservice.NewRepository(db, logger)
	bookService := bookservice.NewService(&bookRepository, logger)
	h := bookservice.NewHandler(&bookService, logger)

	router := http.NewServeMux()
	router = bookservice.InitializeRoutes(router, h)
	router.HandleFunc("GET /api/v1/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    conf.App.Port,
		Handler: middleware.Wraps(router, logger),
	}
	logger.Info("Listening...")

	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
	}
}
