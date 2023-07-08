import { HttpClient } from '@angular/common/http';
import { Injectable, signal, effect } from '@angular/core';
import { tap } from 'rxjs';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly SESSION_NAME = 'loggedIn';
  loggedIn = signal<boolean>(
    sessionStorage.getItem(this.SESSION_NAME) !== null
  );

  constructor(private http: HttpClient) {
    effect(() => {
      if (this.loggedIn()) {
        sessionStorage.setItem(this.SESSION_NAME, 'true');
      } else {
        sessionStorage.removeItem(this.SESSION_NAME);
      }
    });
  }

  login({
    username,
    password,
  }: Partial<{ username: string; password: string }>) {
    return this.http
      .post<{ sessionId: string }>(
        `${environment.baseUrl}/auth/login`,
        {
          username,
          password,
        },
        { withCredentials: true }
      )
      .pipe(
        tap((response) => {
          if (response?.sessionId?.length > 0) {
            this.loggedIn.set(true);
          }
        })
      );
  }

  isLoggedIn() {
    return this.loggedIn.asReadonly();
  }

  logout() {
    return this.http
      .post(`${environment.baseUrl}/auth/logout`, {}, { withCredentials: true })
      .pipe(
        tap(() => {
          this.loggedIn.set(false);
        })
      );
  }
}
