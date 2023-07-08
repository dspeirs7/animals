import {
  Component,
  DestroyRef,
  OnInit,
  Signal,
  ViewChild,
  inject,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule, HttpEventType } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import {
  FormArray,
  FormControl,
  FormGroup,
  ReactiveFormsModule,
} from '@angular/forms';
import { MatNativeDateModule } from '@angular/material/core';
import { MatInputModule } from '@angular/material/input';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { switchMap } from 'rxjs';
import { AuthService } from '../../login/auth.service';
import { Animal, Vaccination } from 'src/app/models/animal';
import { AnimalService } from '../animal.service';
import { AnimalCardComponent } from '../animal-card/animal-card.component';

interface VaccinationForm {
  name: FormControl<string>;
  dateGiven: FormControl<Date>;
  dateNeeded: FormControl<Date>;
}

@Component({
  selector: 'app-animal',
  standalone: true,
  imports: [
    CommonModule,
    HttpClientModule,
    ReactiveFormsModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatTableModule,
    MatSortModule,
    MatProgressBarModule,
    AnimalCardComponent,
  ],
  templateUrl: './animal.component.html',
  styleUrls: ['./animal.component.scss'],
})
export class AnimalComponent implements OnInit {
  @ViewChild(MatSort) sort: MatSort;
  isLoggedIn: Signal<boolean>;
  displayedColumns: string[] = ['name', 'dateGiven', 'dateNeeded', 'delete'];
  vaccinationsForm: FormGroup;
  uploadProgress: number;
  animal = signal<Animal>(undefined as never as Animal);
  today = new Date();
  destroyRef = inject(DestroyRef);

  constructor(
    private authService: AuthService,
    private animalService: AnimalService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.isLoggedIn = this.authService.isLoggedIn();

    this.route.paramMap
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        switchMap((params) => {
          const animalId = params.get('animalId');

          return this.animalService.getAnimal(animalId || '');
        })
      )
      .subscribe((animal) => {
        this.animal.set(animal);
      });

    this.vaccinationsForm = new FormGroup({
      vaccinations: new FormArray([]),
    });
  }

  get vaccinations() {
    return this.vaccinationsForm.get('vaccinations') as FormArray;
  }

  onFileSelected(id: string, event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files as FileList;

    if (files.length) {
      this.animalService
        .addImage(id, files[0])
        .pipe(takeUntilDestroyed(this.destroyRef))
        .subscribe((event) => {
          if (event.type === HttpEventType.UploadProgress) {
            this.uploadProgress = Math.round(
              100 * (event.loaded / (event.total || 1))
            );
          }

          if (event.type === HttpEventType.Response) {
            this.uploadProgress = 0;
            this.animal.mutate((animal) => {
              animal.imageUrl = (event.body as Partial<Animal>).imageUrl!;
            });
          }
        });
    }
  }

  addVaccination() {
    this.vaccinations.push(
      new FormGroup<VaccinationForm>({
        name: new FormControl(),
        dateGiven: new FormControl(),
        dateNeeded: new FormControl(),
      })
    );
  }

  removeVaccination(index: number) {
    this.vaccinations.removeAt(index);
  }

  deleteVaccination(animal: Animal, vaccinationToDelete: Vaccination) {
    this.animalService
      .deleteVaccination(animal.id, vaccinationToDelete)
      .subscribe(() => {
        this.animal.mutate((animal) => {
          animal.vaccinations = animal.vaccinations.filter(
            (vaccination) =>
              vaccination.name !== vaccinationToDelete.name ||
              vaccination.dateGiven !== vaccinationToDelete.dateGiven ||
              vaccination.dateNeeded !== vaccinationToDelete.dateNeeded
          );
        });
      });
  }

  addVaccinations(animal: Animal) {
    this.animalService
      .addVaccinations(animal.id, this.vaccinations.value)
      .subscribe(() => {
        this.animal.mutate((animal) => {
          if (animal.vaccinations?.length) {
            animal.vaccinations.push(this.vaccinations.value);
          } else {
            animal.vaccinations = [this.vaccinations.value];
          }
          this.vaccinations.clear();
        });
      });
  }

  getDataSource(vaccinations: Vaccination[]) {
    const dataSource = new MatTableDataSource(vaccinations);
    dataSource.sort = this.sort;

    return dataSource;
  }
}
