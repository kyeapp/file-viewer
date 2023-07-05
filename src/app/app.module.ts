import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { HttpClientModule } from '@angular/common/http';
import { FolderFileListComponent } from './folder-file-list/folder-file-list.component';
import { SidebarComponent } from './sidebar/sidebar.component';

@NgModule({
  declarations: [
    AppComponent,
    FolderFileListComponent,
    SidebarComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
