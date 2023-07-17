package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dspeirs7/animals/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (a *api) getCats(w http.ResponseWriter, r *http.Request) {
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
}

func (a *api) getChickens(w http.ResponseWriter, r *http.Request) {
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
}

func (a *api) getDogs(w http.ResponseWriter, r *http.Request) {
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
}

func (a *api) getAnimal(w http.ResponseWriter, r *http.Request) {
	animal := r.Context().Value("animal").(*domain.Animal)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

func (a *api) insertAnimal(w http.ResponseWriter, r *http.Request) {
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
}

func (a *api) updateAnimal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

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
}

func (a *api) deleteAnimal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	animal := r.Context().Value("animal").(*domain.Animal)

	if err := a.animalRepo.Delete(ctx, id); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if animal.ImageUrl != "" {
		_ = os.Remove(fmt.Sprintf("./%s", animal.ImageUrl))
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

	if err := a.animalRepo.AddVaccinations(ctx, id, vaccinations); err != nil {
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

	if err := a.animalRepo.DeleteVaccination(ctx, id, vaccination); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (a *api) uploadImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := chi.URLParam(r, "id")

	r.ParseMultipartForm(10 << 20)

	animal := r.Context().Value("animal").(*domain.Animal)
	if animal.ImageUrl != "" {
		go func(imageUrl string) {
			_ = os.Remove(fmt.Sprintf("./%s", imageUrl))
		}(animal.ImageUrl)
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

		id := chi.URLParam(r, "id")

		animal, err := a.animalRepo.GetById(ctx, id)
		if err != nil {
			a.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		ctx = context.WithValue(r.Context(), "animal", animal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
