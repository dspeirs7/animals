import { Component, OnInit } from '@angular/core';
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
import { MatDialogModule } from '@angular/material/dialog';

interface AddChickenForm {
  name: FormControl<string>;
  description: FormControl<string>;
  type: FormControl<number>;
}

@Component({
  selector: 'app-add-chicken-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    ReactiveFormsModule,
  ],
  templateUrl: './add-chicken-dialog.component.html',
  styleUrls: ['./add-chicken-dialog.component.scss'],
})
export class AddChickenDialogComponent implements OnInit {
  addChickenForm: FormGroup;

  ngOnInit(): void {
    this.addChickenForm = new FormGroup<AddChickenForm>({
      name: new FormControl<string>('', { nonNullable: true }),
      description: new FormControl(),
      type: new FormControl<number>(0, {
        nonNullable: true,
        validators: Validators.min(1),
      }),
    });
  }
}
