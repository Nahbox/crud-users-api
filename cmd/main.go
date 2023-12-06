package main

import (
	"context"
	"fmt"
	_ "github.com/Nahbox/crud-users-api/docs"
	"github.com/Nahbox/crud-users-api/internal/config"
	"github.com/Nahbox/crud-users-api/internal/db"
	"github.com/Nahbox/crud-users-api/internal/handler"
	"github.com/Nahbox/crud-users-api/internal/handler/user"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Crud users API
// @version 0.0.1
// @description This API implements CRUD operations, providing versatile access to user data stored in a PostgreSQL database. Users can retrieve a list of users in various formats such as JSON, XML, and TOML. Additionally, the API supports the ability to delete and create individual user records, offering a comprehensive set of functionalities for managing user data seamlessly.

// @host localhost:8080
// @BasePath /
func main() {
	if os.Getenv("APP_ENV") == "local" {
		godotenv.Load()
	}

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal("read config from env", err)
	}

	database, err := db.Initialize(cfg.PgConfig)
	if err != nil {
		log.Fatal("init db", err)
	}
	defer database.Close()

	usersHandler := user.NewHandler(database)

	router := chi.NewRouter()
	router.MethodNotAllowed(handler.MethodNotAllowed)
	router.NotFound(handler.NotFound)

	router.Route("/users", func(router chi.Router) {
		router.Get("/", usersHandler.GetAllUsers)
		router.Post("/", usersHandler.CreateUser)
		router.Route("/{userId}", func(router chi.Router) {
			router.Use(usersHandler.UserContext)
			router.Get("/", usersHandler.GetUser)
			router.Put("/", usersHandler.UpdateUser)
			router.Delete("/", usersHandler.DeleteUser)
		})
	})

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithError(err).Fatal("run http server")
		}
	}()
	defer Stop(server)

	log.Infof("started API server on %s", addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	log.Infoln("stopping API server")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("shutdown server")
		os.Exit(1)
	}
}
