import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: 'chickens' },
  {
    path: 'login',
    loadComponent: () =>
      import('./login/login.component').then((mod) => mod.LoginComponent),
  },
  {
    path: 'chickens',
    loadComponent: () =>
      import('./chickens/chickens.component').then(
        (mod) => mod.ChickensComponent
      ),
  },
  {
    path: 'chickens/:chickenId',
    loadComponent: () =>
      import('./chicken/chicken.component').then((mod) => mod.ChickenComponent),
  },
  {
    path: '**',
    redirectTo: 'chickens',
  },
];
