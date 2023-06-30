import {
  Component,
  DestroyRef,
  OnInit,
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
import { Observable, switchMap } from 'rxjs';
import { Chicken, Vaccination } from '../chicken';
import { ChickenService } from '../chicken.service';
import { ChickenCardComponent } from '../chicken-card/chicken-card.component';
import { AuthService } from '../login/auth.service';

interface VaccinationForm {
  name: FormControl<string>;
  dateGiven: FormControl<Date>;
  dateNeeded: FormControl<Date>;
}

@Component({
  selector: 'app-chicken',
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
    ChickenCardComponent,
  ],
  templateUrl: './chicken.component.html',
  styleUrls: ['./chicken.component.scss'],
})
export class ChickenComponent implements OnInit {
  @ViewChild(MatSort) sort: MatSort;
  isLoggedIn$: Observable<boolean>;
  displayedColumns: string[] = ['name', 'dateGiven', 'dateNeeded', 'delete'];
  vaccinationsForm: FormGroup;
  uploadProgress: number;
  chicken = signal<Chicken>(undefined as never as Chicken);
  today = new Date();
  destroyRef = inject(DestroyRef);

  constructor(
    private authService: AuthService,
    private chickenService: ChickenService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.isLoggedIn$ = this.authService.isLoggedIn();

    this.route.paramMap
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        switchMap((params) => {
          const chickenId = params.get('chickenId');

          return this.chickenService.getChicken(chickenId || '');
        })
      )
      .subscribe((chicken) => {
        this.chicken = signal(chicken);
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
      this.chickenService
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
            this.chicken.set({
              ...this.chicken(),
              imageUrl: (event.body as Partial<Chicken>).imageUrl!,
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

  deleteVaccination(chicken: Chicken, vaccinationToDelete: Vaccination) {
    const updatedChicken = {
      ...chicken,
      vaccinations: chicken.vaccinations.filter(
        (vaccination) => vaccination !== vaccinationToDelete
      ),
    };

    this.chickenService
      .deleteVaccination(chicken.id, vaccinationToDelete)
      .subscribe(() => {
        const updatedChicken = {
          ...chicken,
          vaccinations: chicken.vaccinations.filter(
            (vaccination) => vaccination !== vaccinationToDelete
          ),
        };
        this.chicken.set(updatedChicken);
      });
  }

  addVaccinations(chicken: Chicken) {
    this.chickenService
      .addVaccinations(chicken.id, this.vaccinations.value)
      .subscribe(() => {
        const updatedChicken = {
          ...chicken,
          vaccinations: [
            ...(chicken.vaccinations || []),
            ...this.vaccinations.value,
          ],
        };
        this.chicken.set(updatedChicken);
        this.vaccinations.clear();
      });
  }

  getDataSource(vaccinations: Vaccination[]) {
    const dataSource = new MatTableDataSource(vaccinations);
    dataSource.sort = this.sort;

    return dataSource;
  }
}
