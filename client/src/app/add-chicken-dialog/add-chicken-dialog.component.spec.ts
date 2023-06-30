import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddChickenDialogComponent } from './add-chicken-dialog.component';

describe('AddChickenDialogComponent', () => {
  let component: AddChickenDialogComponent;
  let fixture: ComponentFixture<AddChickenDialogComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AddChickenDialogComponent]
    });
    fixture = TestBed.createComponent(AddChickenDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
