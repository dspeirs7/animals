package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	"golang.org/x/crypto/bcrypt"
)

type api struct {
	logger   *zap.Logger
	dbClient *mongo.Client

	chickenRepo domain.ChickenRepository
}

var sessions = make(map[string]domain.Session)

func NewAPI(ctx context.Context, logger *zap.Logger) *api {
	dbClient := repository.GetDB(ctx)
	db := dbClient.Database("chickens")

	chickenRepo := repository.NewChickenRepository(db.Collection("Chickens"), db.Collection("users"))

	return &api{
		logger:   logger,
		dbClient: dbClient,

		chickenRepo: chickenRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *chi.Mux {
	env := os.Getenv("ENV")

	r := chi.NewRouter()

	r.Use(chim.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiRouter := chi.NewRouter()

	apiRouter.Post("/login", a.login)
	apiRouter.Post("/logout", a.logout)

	apiRouter.Route("/image/{id}", func(r chi.Router) {
		r.Use(middleware.GetSession(sessions))
		r.Use(a.ChickenCtx)
		r.Post("/", a.uploadImage)
	})

	apiRouter.Route("/chickens", func(r chi.Router) {
		r.Get("/", a.getChickens)
		r.With(middleware.GetSession(sessions)).Post("/", a.insertChicken)
	})
	apiRouter.Route("/chicken/{id}", func(r chi.Router) {
		r.Use(a.ChickenCtx)
		r.Get("/", a.getChicken)
		r.Group(func(r chi.Router) {
			r.Use(middleware.GetSession(sessions))
			r.Put("/", a.updateChicken)
			r.Delete("/", a.deleteChicken)
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

func (a *api) login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var user domain.User

	if err := decoder.Decode(&user); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	admin, err := a.chickenRepo.GetUser(ctx, "admin")
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(user.Password)); err != nil {
		a.errorResponse(w, r, http.StatusUnauthorized, err)
		return
	}

	sessionId := domain.SessionId()
	expiresAt := time.Now().Add(3600 * time.Second)
	sessions[sessionId] = domain.Session{Username: admin.Username, Expiry: expiresAt}

	http.SetCookie(w, &http.Cookie{Name: "session_token", Value: sessionId, Expires: expiresAt, Path: "/", SameSite: http.SameSiteLaxMode})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"sessionId": sessionId})
}

func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		a.logger.Info(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionId := cookie.Value
	delete(sessions, sessionId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (a *api) getChickens(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	results, err := a.chickenRepo.GetAll(ctx)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (a *api) getChicken(w http.ResponseWriter, r *http.Request) {
	chicken := r.Context().Value("chicken").(*domain.Chicken)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chicken)
}

func (a *api) insertChicken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var chicken domain.Chicken

	if err := decoder.Decode(&chicken); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	result, err := a.chickenRepo.Insert(ctx, chicken)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (a *api) updateChicken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	decoder := json.NewDecoder(r.Body)
	var chicken domain.Chicken

	if err := decoder.Decode(&chicken); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := a.chickenRepo.Update(ctx, id, chicken); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteChicken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	chicken := r.Context().Value("chicken").(*domain.Chicken)

	if err := a.chickenRepo.Delete(ctx, id); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if chicken.ImageUrl != "" {
		_ = os.Remove(fmt.Sprintf("./%s", chicken.ImageUrl))
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) addVaccinations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	decoder := json.NewDecoder(r.Body)

	var vaccinations []domain.Vaccination

	if err := decoder.Decode(&vaccinations); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := a.chickenRepo.AddVaccinations(ctx, id, vaccinations); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (a *api) deleteVaccination(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	decoder := json.NewDecoder(r.Body)

	var vaccination domain.Vaccination

	if err := decoder.Decode(&vaccination); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := a.chickenRepo.DeleteVaccination(ctx, id, vaccination); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (a *api) uploadImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	r.ParseMultipartForm(10 << 20)

	chicken := r.Context().Value("chicken").(*domain.Chicken)

	if chicken.ImageUrl != "" {
		_ = os.Remove(fmt.Sprintf("./%s", chicken.ImageUrl))
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	defer file.Close()

	if err := os.MkdirAll(filepath.Join(".", "images"), os.ModePerm); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	fileName := fmt.Sprintf("images/%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))

	dst, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
	}

	if err := a.chickenRepo.UpdateImageUrl(ctx, id, fileName); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.Chicken{ImageUrl: fileName})
}

func (a *api) ChickenCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := chi.URLParam(r, "id")

		chicken, err := a.chickenRepo.GetById(ctx, id)
		if err != nil {
			a.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		ctx = context.WithValue(r.Context(), "chicken", chicken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
