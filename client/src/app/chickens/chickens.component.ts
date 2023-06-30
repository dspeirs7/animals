import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Observable, filter, switchMap } from 'rxjs';
import { HttpClientModule } from '@angular/common/http';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { Router, RouterModule } from '@angular/router';
import { ChickenService } from '../chicken.service';
import { AddChickenDialogComponent } from '../add-chicken-dialog/add-chicken-dialog.component';
import { Chicken } from '../chicken';
import { ChickenCardComponent } from '../chicken-card/chicken-card.component';
import { AuthService } from '../login/auth.service';

@Component({
  selector: 'app-chickens',
  standalone: true,
  imports: [
    CommonModule,
    HttpClientModule,
    RouterModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    ChickenCardComponent,
  ],
  templateUrl: './chickens.component.html',
  styleUrls: ['./chickens.component.scss'],
})
export class ChickensComponent {
  isLoggedIn$: Observable<boolean>;
  chickens = signal<Chicken[]>([]);

  constructor(
    private authService: AuthService,
    private chickenService: ChickenService,
    private matDialog: MatDialog,
    private router: Router
  ) {}

  ngOnInit() {
    this.isLoggedIn$ = this.authService.isLoggedIn();

    this.chickenService
      .getChickens()
      .subscribe((chickens) => this.chickens.set(chickens));
  }

  addChicken() {
    const dialogRef = this.matDialog.open(AddChickenDialogComponent);

    dialogRef
      .afterClosed()
      .pipe(
        filter((chicken) => chicken),
        switchMap((chicken) => this.chickenService.addChicken(chicken))
      )
      .subscribe((chicken) => {
        this.router.navigate(['/', 'chickens', chicken.id]);
      });
  }

  onDelete(id: string) {
    this.chickens.set([
      ...this.chickens().filter((chicken) => chicken.id != id),
    ]);
  }
}
