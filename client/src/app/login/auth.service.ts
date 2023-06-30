import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private loggedInSubject = new BehaviorSubject<boolean>(false);

  constructor(private http: HttpClient) {}

  login(username: string, password: string) {
    return this.http
      .post<{ sessionId: string }>(
        `${environment.apiUrl}/login`,
        {
          username,
          password,
        },
        { withCredentials: true }
      )
      .pipe(
        tap((response) => {
          if (response?.sessionId?.length > 0) {
            this.loggedInSubject.next(true);
          }
        })
      );
  }

  isLoggedIn() {
    return this.loggedInSubject.asObservable();
  }

  logout() {
    return this.http
      .post(`${environment.apiUrl}/logout`, {}, { withCredentials: true })
      .pipe(
        tap(() => {
          this.loggedInSubject.next(false);
        })
      );
  }
}
