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
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';

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
export class AppComponent {
  location = inject(Location);
  router = inject(Router);
  authService = inject(AuthService);
  isLoggedIn = this.authService.isLoggedIn();
  destroyRef = inject(DestroyRef);
  showBack = toSignal(
    this.router.events.pipe(
      takeUntilDestroyed(this.destroyRef),
      filter((event) => event instanceof NavigationEnd),
      map(
        (event) =>
          !['/cats', '/chickens', '/dogs'].includes(
            (event as NavigationEnd).urlAfterRedirects
          )
      )
    )
  );

  logout() {
    this.authService.logout().subscribe(() => {
      this.router.navigate(['/']);
    });
  }

  goBack() {
    this.location.back();
  }
}
