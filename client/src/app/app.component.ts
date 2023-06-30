import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  RouterOutlet,
  RouterModule,
  Router,
  NavigationEnd,
} from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { Observable, filter, map } from 'rxjs';
import { AuthService } from './login/auth.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    RouterModule,
    MatButtonModule,
    MatIconModule,
    MatToolbarModule,
  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  url$: Observable<string>;
  isLoggedIn$: Observable<boolean>;

  constructor(private router: Router, private authService: AuthService) {}

  ngOnInit(): void {
    this.url$ = this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map((event) => (event as NavigationEnd).url)
    );

    this.isLoggedIn$ = this.authService.isLoggedIn();
  }

  logout() {
    this.authService.logout().subscribe(() => {
      this.router.navigate(['/chickens']);
    });
  }
}
