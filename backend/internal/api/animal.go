package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dspeirs7/animals/internal/domain"
)

func (a *api) getCats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		results, err := a.animalRepo.GetAllCats(ctx)
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) getChickens(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		results, err := a.animalRepo.GetAllChickens(ctx)
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) getDogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		results, err := a.animalRepo.GetAllDogs(ctx)
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) handleAnimal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := path.Base(r.URL.Path)
	animal := r.Context().Value("animal").(*domain.Animal)

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(animal)
	case http.MethodPost:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		decoder := json.NewDecoder(r.Body)
		var animal domain.Animal

		if err := decoder.Decode(&animal); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		result, err := a.animalRepo.Insert(ctx, animal)
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var animal domain.Animal

		if err := decoder.Decode(&animal); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := a.animalRepo.Update(ctx, id, animal); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		if err := a.animalRepo.Delete(ctx, id); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if animal.ImageUrl != "" {
			_ = os.Remove(fmt.Sprintf("./%s", animal.ImageUrl))
		}

		w.WriteHeader(http.StatusOK)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, PUT, POST, DELETE, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, PUT, POST, DELETE, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) addVaccinations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		id := path.Base(r.URL.Path)

		decoder := json.NewDecoder(r.Body)

		var vaccinations []domain.Vaccination

		if err := decoder.Decode(&vaccinations); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := a.animalRepo.AddVaccinations(ctx, id, vaccinations); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) deleteVaccination(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		id := path.Base(r.URL.Path)

		decoder := json.NewDecoder(r.Body)

		var vaccination domain.Vaccination

		if err := decoder.Decode(&vaccination); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := a.animalRepo.DeleteVaccination(ctx, id, vaccination); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) uploadImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := path.Base(r.URL.Path)

	r.ParseMultipartForm(10 << 20)

	animal := r.Context().Value("animal").(*domain.Animal)

	if animal.ImageUrl != "" {
		_ = os.Remove(fmt.Sprintf("./%s", animal.ImageUrl))
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

	if err := a.animalRepo.UpdateImageUrl(ctx, id, fileName); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.Animal{ImageUrl: fileName})
}

func (a *api) AnimalCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := path.Base(r.URL.Path)

		animal, err := a.animalRepo.GetById(ctx, id)
		if err != nil {
			a.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		ctx = context.WithValue(r.Context(), "animal", animal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
