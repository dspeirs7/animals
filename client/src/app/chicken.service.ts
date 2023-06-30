import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Chicken, Vaccination } from './chicken';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ChickenService {
  constructor(private http: HttpClient) {}

  getChickens() {
    return this.http.get<Chicken[]>(`${environment.apiUrl}/chickens`);
  }

  getChicken(id: string) {
    return this.http.get<Chicken>(`${environment.apiUrl}/chicken/${id}`);
  }

  addChicken(chicken: Partial<Chicken>): Observable<Chicken> {
    return this.http.post<Chicken>(
      `${environment.apiUrl}/chickens`,
      {
        ...chicken,
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

  updateChicken(chicken: Chicken) {
    return this.http.put(
      `${environment.apiUrl}/chicken/${chicken.id}`,
      chicken,
      { withCredentials: true }
    );
  }

  addVaccinations(id: string, vaccinations: Vaccination[]) {
    return this.http.post(
      `${environment.apiUrl}/chicken/${id}/vaccinations/add`,
      vaccinations,
      { withCredentials: true }
    );
  }

  deleteVaccination(id: string, vaccination: Vaccination) {
    return this.http.post(
      `${environment.apiUrl}/chicken/${id}/vaccinations/delete`,
      vaccination,
      { withCredentials: true }
    );
  }

  deleteChicken(id: string) {
    return this.http.delete(`${environment.apiUrl}/chicken/${id}`, {
      withCredentials: true,
    });
  }
}
