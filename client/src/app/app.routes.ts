import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: 'cats' },
  {
    path: 'login',
    loadComponent: () =>
      import('./login/login.component').then((mod) => mod.LoginComponent),
  },
  {
    path: 'cats',
    data: { animalType: 1 },
    loadComponent: () =>
      import('./animals/animals.component').then((mod) => mod.AnimalsComponent),
  },
  {
    path: 'chickens',
    data: { animalType: 2 },
    loadComponent: () =>
      import('./animals/animals.component').then((mod) => mod.AnimalsComponent),
  },
  {
    path: 'dogs',
    data: { animalType: 3 },
    loadComponent: () =>
      import('./animals/animals.component').then((mod) => mod.AnimalsComponent),
  },
  {
    path: 'animal/:animalId',
    loadComponent: () =>
      import('./animals/animal/animal.component').then(
        (mod) => mod.AnimalComponent
      ),
  },
  {
    path: '**',
    redirectTo: 'cats',
  },
];
