import { Component, Inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormControl,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';

interface AddChickenForm {
  name: FormControl<string>;
  description: FormControl<string>;
  type: FormControl<number>;
  breed: FormControl<number>;
}

interface DialogData {
  animalType: 1 | 2 | 3;
}

@Component({
  selector: 'app-add-animal-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    ReactiveFormsModule,
  ],
  templateUrl: './add-animal-dialog.component.html',
  styleUrls: ['./add-animal-dialog.component.scss'],
})
export class AddAnimalDialogComponent {
  constructor(@Inject(MAT_DIALOG_DATA) public data: DialogData) {}

  addAnimalForm: FormGroup;

  ngOnInit(): void {
    this.addAnimalForm = new FormGroup<AddChickenForm>({
      name: new FormControl<string>('', { nonNullable: true }),
      description: new FormControl(),
      type: new FormControl<number>(this.data.animalType, {
        nonNullable: true,
      }),
      breed: new FormControl<number>(0, {
        nonNullable: true,
        validators: Validators.min(1),
      }),
    });
  }
}
