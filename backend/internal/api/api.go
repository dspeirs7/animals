package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dspeirs7/animals/internal/domain"
	"github.com/dspeirs7/animals/internal/middleware"
	"github.com/dspeirs7/animals/internal/repository"
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      a.Routes(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (a *api) Routes() *chi.Mux {
	env := os.Getenv("ENV")

	r := chi.NewRouter()

	r.Use(chim.Logger)

	if env != "production" && env != "dev" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:4200"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	r.Post("/auth/login", a.login)
	r.Post("/auth/logout", a.logout)

	apiRouter := chi.NewRouter()

	apiRouter.Route("/image/{id}", func(r chi.Router) {
		r.Use(middleware.GetSession)
		r.Use(a.AnimalCtx)
		r.Post("/", a.uploadImage)
	})

	apiRouter.Get("/cats", a.getCats)

	apiRouter.Get("/chickens", a.getChickens)

	apiRouter.Get("/dogs", a.getDogs)

	apiRouter.Route("/animal", func(r chi.Router) {
		r.Use(middleware.GetSession)
		r.Post("/", a.insertAnimal)
	})

	apiRouter.Route("/animal/{id}", func(r chi.Router) {
		r.Use(a.AnimalCtx)
		r.Get("/", a.getAnimal)
		r.Group(func(r chi.Router) {
			r.Use(middleware.GetSession)
			r.Put("/", a.updateAnimal)
			r.Delete("/", a.deleteAnimal)
			r.Post("/vaccinations/add", a.addVaccinations)
			r.Post("/vaccinations/delete", a.deleteVaccination)
		})
	})

	r.Mount("/api", apiRouter)

	fs := http.FileServer(http.Dir("images"))
	r.Handle("/images/*", http.StripPrefix("/images/", fs))

	if env == "production" || env == "dev" {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			workDir, _ := os.Getwd()
			filesDir := filepath.Join(workDir, "dist")

			if _, err := os.Stat(filesDir + r.URL.Path); errors.Is(err, os.ErrNotExist) {
				http.ServeFile(w, r, filepath.Join(filesDir, "index.html"))
			}
			http.ServeFile(w, r, filesDir+r.URL.Path)
		})
	}

	return r
}

func (a *api) Disconnect(ctx context.Context) error {
	return a.dbClient.Disconnect(ctx)
}
