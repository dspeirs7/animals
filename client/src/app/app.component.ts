import {
  Component,
  DestroyRef,
  OnInit,
  Signal,
  inject,
  signal,
} from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import {
  RouterOutlet,
  RouterModule,
  Router,
  NavigationEnd,
} from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { filter, map } from 'rxjs';
import { AuthService } from './login/auth.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

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
  isLoggedIn: Signal<boolean>;
  showBack = signal<boolean>(false);
  destroyRef = inject(DestroyRef);

  constructor(
    private location: Location,
    private router: Router,
    private authService: AuthService
  ) {
    this.router.events
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        filter((event) => event instanceof NavigationEnd)
      )
      .subscribe((event) => {
        this.showBack.set(
          !['/cats', '/chickens', '/dogs'].includes(
            (event as NavigationEnd).urlAfterRedirects
          )
        );
      });
  }

  ngOnInit(): void {
    this.isLoggedIn = this.authService.isLoggedIn();
  }

  logout() {
    this.authService.logout().subscribe(() => {
      this.router.navigate(['/']);
    });
  }

  goBack() {
    this.location.back();
  }
}
