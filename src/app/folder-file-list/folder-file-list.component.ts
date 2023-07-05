import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'folder-file-list',
  templateUrl: './folder-file-list.component.html',
  styleUrls: ['./folder-file-list.component.css']
})
export class FolderFileListComponent {
  dirEntries: FsEntry[] = [];
  selectedEntry: FsEntry | null = null;

  constructor(private http: HttpClient) {}

  selectEntry(entry: FsEntry) {
    this.selectedEntry = entry;
  }

  async ngOnInit() {
    this.dirEntries = await this.fetchDirectoryEntries();
  }

  fetchDirectoryEntries(): Promise<FsEntry[]> {
    const url = 'http://localhost:8080/directory';
    return this.http.get<FsEntry[]>(url).toPromise()
    .then(response => response || []);
  }
}

interface FsEntry {
  name: string;
  isFile: boolean;
  isFolder: boolean;
}
