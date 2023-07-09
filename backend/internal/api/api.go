package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dspeirs7/animals/internal/domain"
	"github.com/dspeirs7/animals/internal/middleware"
	"github.com/dspeirs7/animals/internal/repository"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type api struct {
	logger   *zap.Logger
	dbClient *mongo.Client

	animalRepo domain.AnimalRepository
	userRepo   domain.UserRepository
}

var sessions = make(map[string]domain.Session)

func NewAPI(ctx context.Context, logger *zap.Logger) *api {
	dbClient := repository.GetDB(ctx)
	db := dbClient.Database("animals")

	animalRepo := repository.NewAnimalRepository(db.Collection("animals"))
	userRepo := repository.NewUserRepository(db.Collection("users"))

	return &api{
		logger:   logger,
		dbClient: dbClient,

		animalRepo: animalRepo,
		userRepo:   userRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	env := os.Getenv("ENV")

	var handler http.Handler

	if env != "production" && env == "dev" {
		handler = cors.Default().Handler(a.Routes())
	} else {
		handler = a.Routes()
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}

func (a *api) Routes() *http.ServeMux {
	env := os.Getenv("ENV")

	r := http.NewServeMux()

	r.HandleFunc("/auth/login", a.login)
	r.HandleFunc("/auth/logout", a.logout)

	r.Handle("/api/image/", middleware.Logger(middleware.GetSession(sessions)(a.AnimalCtx(http.HandlerFunc(a.uploadImage)))))
	r.Handle("/api/cats", middleware.Logger(http.HandlerFunc(a.getCats)))
	r.Handle("/api/chickens", middleware.Logger(http.HandlerFunc(a.getChickens)))
	r.Handle("/api/dogs", middleware.Logger(http.HandlerFunc(a.getDogs)))
	r.Handle("/api/animal/", middleware.Logger(middleware.GetSession(sessions)(a.AnimalCtx(http.HandlerFunc(a.handleAnimal)))))
	r.Handle("/api/vaccination/add/", middleware.Logger(middleware.GetSession(sessions)(http.HandlerFunc(a.addVaccinations))))
	r.Handle("/api/vaccination/delete/", middleware.Logger(middleware.GetSession(sessions)(http.HandlerFunc(a.deleteVaccination))))

	fs := http.FileServer(http.Dir("images"))
	r.Handle("/images/", http.StripPrefix("/images/", fs))

	if env == "production" || env == "dev" {
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			workDir, _ := os.Getwd()
			filesDir := filepath.Join(workDir, "dist")

			if _, err := os.Stat(filesDir + r.URL.Path); errors.Is(err, os.ErrNotExist) {
				http.ServeFile(w, r, filepath.Join(filesDir, "index.html"))
				return
			}
			http.ServeFile(w, r, filesDir+r.URL.Path)
		})
	}

	return r
}

func (a *api) Disconnect(ctx context.Context) error {
	return a.dbClient.Disconnect(ctx)
}
