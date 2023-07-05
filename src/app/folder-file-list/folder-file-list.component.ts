import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'folder-file-list',
  templateUrl: './folder-file-list.component.html',
  styleUrls: ['./folder-file-list.component.css']
})
export class FolderFileListComponent {
  filenames: string[] = [];
  selectedFile: string | null = null;

  constructor(private http: HttpClient) {}

  selectFile(file: string) {
    this.selectedFile = file;
  }

  async ngOnInit() {
    this.filenames = await this.fetchFilenames();
  }

  fetchFilenames(): Promise<string[]> {
    const url = 'http://localhost:8080/files';
    return this.http.get<string[]>(url).toPromise()
    .then(response => response || []);
  }
}
