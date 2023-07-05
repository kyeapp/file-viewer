import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FolderFileListComponent } from './folder-file-list.component';

describe('FolderFileListComponent', () => {
  let component: FolderFileListComponent;
  let fixture: ComponentFixture<FolderFileListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FolderFileListComponent]
    });
    fixture = TestBed.createComponent(FolderFileListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
