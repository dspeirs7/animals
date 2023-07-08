import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { Animal, Vaccination } from '../models/animal';

@Injectable({
  providedIn: 'root',
})
export class AnimalService {
  constructor(private http: HttpClient) {}

  getCats() {
    return this.http.get<Animal[]>(`${environment.apiUrl}/cats`);
  }

  getChickens() {
    return this.http.get<Animal[]>(`${environment.apiUrl}/chickens`);
  }

  getDogs() {
    return this.http.get<Animal[]>(`${environment.apiUrl}/dogs`);
  }

  getAnimal(id: string) {
    return this.http.get<Animal>(`${environment.apiUrl}/animal/${id}`);
  }

  addAnimal(animal: Partial<Animal>): Observable<Animal> {
    return this.http.post<Animal>(
      `${environment.apiUrl}/animal`,
      {
        ...animal,
      },
      { withCredentials: true }
    );
  }

  addImage(id: string, image: File) {
    const formData = new FormData();
    formData.append('image', image, image.name);

    return this.http.post<string>(
      `${environment.apiUrl}/image/${id}`,
      formData,
      {
        reportProgress: true,
        observe: 'events',
        withCredentials: true,
      }
    );
  }

  updateAnimal(animal: Animal) {
    return this.http.put(`${environment.apiUrl}/animal/${animal.id}`, animal, {
      withCredentials: true,
    });
  }

  addVaccinations(id: string, vaccinations: Vaccination[]) {
    return this.http.post(
      `${environment.apiUrl}/animal/${id}/vaccinations/add`,
      vaccinations,
      { withCredentials: true }
    );
  }

  deleteVaccination(id: string, vaccination: Vaccination) {
    return this.http.post(
      `${environment.apiUrl}/animal/${id}/vaccinations/delete`,
      vaccination,
      { withCredentials: true }
    );
  }

  deleteAnimal(id: string) {
    return this.http.delete(`${environment.apiUrl}/animal/${id}`, {
      withCredentials: true,
    });
  }
}
