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
	"github.com/dspeirs7/animals/internal/repository"
	"github.com/dspeirs7/mongostore"
	"github.com/gorilla/mux"
	secrets "github.com/ijustfool/docker-secrets"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type api struct {
	logger   *zap.Logger
	dbClient *mongo.Client

	animalRepo domain.AnimalRepository
	userRepo   domain.UserRepository
	store      *mongostore.MongoStore
}

func NewAPI(ctx context.Context, logger *zap.Logger) *api {
	var sessionKey string

	dockerSecrets, _ := secrets.NewDockerSecrets("")

	sessionKey, err := dockerSecrets.Get("secret_key")
	if err != nil {
		sessionKey = os.Getenv("SECRET_KEY")
	}

	dbClient := repository.GetDB(ctx)
	db := dbClient.Database("animals")

	animalRepo := repository.NewAnimalRepository(db.Collection("animals"))
	userRepo := repository.NewUserRepository(db.Collection("users"))
	store := mongostore.NewMongoStore(db.Collection("sessions"), 3600, true, []byte(sessionKey))

	return &api{
		logger:   logger,
		dbClient: dbClient,

		animalRepo: animalRepo,
		userRepo:   userRepo,
		store:      store,
	}
}

func (a *api) Server(port int) *http.Server {
	env := os.Getenv("ENV")

	var handler http.Handler

	if env != "production" && env != "dev" {
		handler = cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:4200"},
			AllowCredentials: true,
		}).Handler(a.Routes())
	} else {
		handler = a.Routes()
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (a *api) Routes() *mux.Router {
	env := os.Getenv("ENV")

	r := mux.NewRouter()

	r.Use(a.Logger)

	r.HandleFunc("/auth/login", a.login).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", a.logout).Methods(http.MethodPost)

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/cats", a.getCats).Methods(http.MethodGet)
	apiRouter.HandleFunc("/chickens", a.getChickens).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dogs", a.getDogs).Methods(http.MethodGet)

	imageRouter := apiRouter.PathPrefix("/image/{id}").Subrouter()
	imageRouter.HandleFunc("/", a.uploadImage).Methods(http.MethodPost)
	imageRouter.Use(a.GetSession)
	imageRouter.Use(a.AnimalCtx)

	animalRouter := apiRouter.PathPrefix("/animal").Subrouter()

	animalRouter.HandleFunc("/", a.insertAnimal).Methods(http.MethodPost)
	animalRouter.HandleFunc("/{id}", a.getAnimal).Methods(http.MethodGet)
	animalRouter.HandleFunc("/{id}", a.updateAnimal).Methods(http.MethodPut)
	animalRouter.HandleFunc("/{id}", a.deleteAnimal).Methods(http.MethodDelete)
	animalRouter.HandleFunc("/{id}/vaccinations/add", a.addVaccinations).Methods(http.MethodPost)
	animalRouter.HandleFunc("/{id}/vaccinations/delete", a.deleteVaccination).Methods(http.MethodPost)
	animalRouter.Use(a.GetSession)
	animalRouter.Use(a.AnimalCtx)

	fs := http.FileServer(http.Dir("images"))
	r.Handle("/images", http.StripPrefix("/images/", fs))

	if env == "production" || env == "dev" {
		r.HandleFunc("/", a.handleSpa).Methods(http.MethodGet)
	}

	return r
}

func (a *api) handleSpa(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "dist")

	if _, err := os.Stat(filesDir + r.URL.Path); errors.Is(err, os.ErrNotExist) {
		http.ServeFile(w, r, filepath.Join(filesDir, "index.html"))
	}
	http.ServeFile(w, r, filesDir+r.URL.Path)
}

func (a *api) Disconnect(ctx context.Context) error {
	return a.dbClient.Disconnect(ctx)
}

func (a *api) GetSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			session, err := a.store.Get(r, "session_token")
			if err != nil || session.ID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *api) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := domain.NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		fmt.Printf("%s: \"%s %s\" from %s [%d]-%s\n", time.Now().Format("2006-01-02 15:04:05"),
			r.Method, r.URL.String(), r.RemoteAddr, lrw.StatusCode, http.StatusText(lrw.StatusCode))
	})
}
