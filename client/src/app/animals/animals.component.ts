import {
  Component,
  DestroyRef,
  OnInit,
  Signal,
  computed,
  inject,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { filter, of, switchMap } from 'rxjs';
import { HttpClientModule } from '@angular/common/http';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { AuthService } from '../login/auth.service';
import { AnimalService } from '../animals/animal.service';
import { AnimalCardComponent } from '../animals/animal-card/animal-card.component';
import { Animal, AnimalName, AnimalType } from '../models/animal';
import { AddAnimalDialogComponent } from '../animals/add-animal-dialog/add-animal-dialog.component';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-animals',
  standalone: true,
  imports: [
    CommonModule,
    HttpClientModule,
    RouterModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    AnimalCardComponent,
  ],
  templateUrl: './animals.component.html',
  styleUrls: ['./animals.component.scss'],
})
export class AnimalsComponent implements OnInit {
  route = inject(ActivatedRoute);
  authService = inject(AuthService);
  animalService = inject(AnimalService);
  matDialog = inject(MatDialog);
  router = inject(Router);
  isLoggedIn = this.authService.isLoggedIn();
  animalType = signal<AnimalType>(1);
  animalTypeName = computed<AnimalName>(() => {
    switch (this.animalType()) {
      case 1:
        return 'Cat';
      case 2:
        return 'Chicken';
      case 3:
        return 'Dog';
      default:
        return 'Cat';
    }
  });
  animals = signal<Animal[]>([]);
  destroyRef = inject(DestroyRef);

  ngOnInit(): void {
    this.route.data
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        switchMap((data) => {
          const animalType = data['animalType'];
          this.animalType.set(animalType);

          switch (animalType) {
            case 1:
              return this.animalService.getCats();
            case 2:
              return this.animalService.getChickens();
            case 3:
              return this.animalService.getDogs();
          }

          return of([]);
        })
      )
      .subscribe((animals) => {
        this.animals.set(animals || []);
      });
  }

  addAnimal() {
    const dialogRef = this.matDialog.open(AddAnimalDialogComponent, {
      data: { animalType: this.animalType() },
    });

    dialogRef
      .afterClosed()
      .pipe(
        filter((animal) => animal),
        switchMap((animal) => this.animalService.addAnimal(animal))
      )
      .subscribe((animal) => {
        this.router.navigate(['/', 'animal', animal.id]);
      });
  }

  onDelete(id: string) {
    this.animals.set([...this.animals().filter((animal) => animal.id != id)]);
  }
}
