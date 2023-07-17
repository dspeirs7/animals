import {
  Component,
  EventEmitter,
  Input,
  Output,
  Signal,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogModule,
} from '@angular/material/dialog';
import { EnvironmentPipe } from '../../environment.pipe';
import { AuthService } from '../../login/auth.service';
import { Animal, AnimalName } from 'src/app/models/animal';
import { AnimalService } from '../animal.service';

@Component({
  selector: 'app-animal-card',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatDialogModule,
    EnvironmentPipe,
  ],
  templateUrl: './animal-card.component.html',
  styleUrls: ['./animal-card.component.scss'],
})
export class AnimalCardComponent {
  authService = inject(AuthService);
  animalService = inject(AnimalService);
  matDialog = inject(MatDialog);
  @Input() animal: Animal;
  @Input() type: AnimalName;
  @Input() showActions: boolean = false;
  @Output() onDelete = new EventEmitter<string>();
  isLoggedIn: Signal<boolean>;

  ngOnInit(): void {
    this.isLoggedIn = this.authService.isLoggedIn();
  }

  deleteAnimal(animal: Animal) {
    const dialogRef = this.matDialog.open(DeleteAnimalDialog, {
      data: {
        name: animal.name,
        id: animal.id,
      },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.animalService.deleteAnimal(animal.id).subscribe(() => {
          this.onDelete.emit(animal.id);
        });
      }
    });
  }
}

interface DialogData {
  id: string;
  name: string;
}

@Component({
  selector: 'app-delete-animal-dialog',
  standalone: true,
  template: `
    <h2 mat-dialog-title>Confirm Delete</h2>
    <mat-dialog-content>
      Are you sure you want to delete {{ data.name }}?
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-button mat-dialog-close>Cancel</button>
      <button mat-button [mat-dialog-close]="true" cdkFocusInitial>
        Delete
      </button>
    </mat-dialog-actions>
  `,
  imports: [MatDialogModule, MatButtonModule],
})
export class DeleteAnimalDialog {
  data: DialogData = inject(MAT_DIALOG_DATA);
}
