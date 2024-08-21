package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/varonikp/keys-ms/internal/config"
	"github.com/varonikp/keys-ms/internal/repository/pgrepo"
	"github.com/varonikp/keys-ms/internal/services"
	"github.com/varonikp/keys-ms/internal/transport/httpserver"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func run() error {
	conf := config.Read()

	pgDB, err := sqlx.Connect("mysql", conf.DSN)
	if err != nil {
		return fmt.Errorf("sqlx connect failed: %w", err)
	}

	logger := slog.Default()

	logger.Debug("Running MySQL migrations")
	if err := runPostgresMigrations(conf.DSN, conf.MigrationsPath); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	userRepo := pgrepo.NewUserRepository(pgDB)
	softwareRepo := pgrepo.NewSoftwareRepo(pgDB)
	licenseRepo := pgrepo.NewLicenseRepository(pgDB)

	tokenService := services.NewTokenService([]byte("secret"), time.Hour)

	userService := services.NewUserService(userRepo)
	softwareService := services.NewSoftwareService(softwareRepo)
	licenseService := services.NewLicenseService(licenseRepo)

	httpServer := httpserver.NewHttpServer(userService, tokenService, softwareService, licenseService)

	router := mux.NewRouter()

	router.HandleFunc("/signin", httpServer.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/signup", httpServer.SignUp).Methods(http.MethodPost)

	// admin permissions sections
	router.HandleFunc("/admin", httpServer.CheckAdmin(httpServer.GrantAdmin)).Methods(http.MethodPost)
	router.HandleFunc("/admins/{tag}", httpServer.CheckAdmin(httpServer.RevokeAdmin)).Methods(http.MethodDelete)

	// users section
	router.HandleFunc("/users", httpServer.CheckAdmin(httpServer.GetUsers)).Methods(http.MethodGet)
	router.HandleFunc("/user/{tag}", httpServer.CheckAuthorized(httpServer.GetUser)).Methods(http.MethodGet)
	router.HandleFunc("/user/{tag}", httpServer.CheckAdmin(httpServer.UpdateUser)).Methods(http.MethodPatch)
	router.HandleFunc("/user/{tag}", httpServer.CheckAdmin(httpServer.DeleteUser)).Methods(http.MethodDelete)

	// licenses section
	router.HandleFunc("/licenses/{user_id}", httpServer.CheckAuthorized(httpServer.GetLicenses)).Methods(http.MethodGet)
	router.HandleFunc("/license/{license_id}", httpServer.CheckAdmin(httpServer.GetLicense)).Methods(http.MethodGet)
	router.HandleFunc("/license", httpServer.CheckAdmin(httpServer.CreateLicense)).Methods(http.MethodPost)
	router.HandleFunc("/license/{license_id}", httpServer.CheckAdmin(httpServer.UpdateLicense)).Methods(http.MethodPatch)
	router.HandleFunc("/license/{license_id}", httpServer.CheckAdmin(httpServer.DeleteLicense)).Methods(http.MethodDelete)

	// softwares section
	router.HandleFunc("/softwares", httpServer.GetSoftwares).Methods(http.MethodGet)
	router.HandleFunc("/software/{software_id}", httpServer.GetSoftware).Methods(http.MethodGet)
	router.HandleFunc("/software", httpServer.CheckAdmin(httpServer.CreateSoftware)).Methods(http.MethodPost)
	router.HandleFunc("/software/{software_id}", httpServer.CheckAdmin(httpServer.UpdateSoftware)).Methods(http.MethodPatch)
	router.HandleFunc("/software/{software_id}", httpServer.CheckAdmin(httpServer.DeleteSoftware)).Methods(http.MethodDelete)

	srv := &http.Server{
		Addr:    conf.HttpAddr,
		Handler: router,
	}

	slog.Info("server started")

	return srv.ListenAndServe()

}

func runPostgresMigrations(dsn, path string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}

	if dsn == "" {
		return errors.New("no DSN provided")
	}

	m, err := migrate.New(path, dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
