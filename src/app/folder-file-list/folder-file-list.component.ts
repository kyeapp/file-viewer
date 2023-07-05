import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { join, Path } from '@angular-devkit/core';

@Component({
  selector: 'folder-file-list',
  templateUrl: './folder-file-list.component.html',
  styleUrls: ['./folder-file-list.component.css']
})
export class FolderFileListComponent {
  constructor(private http: HttpClient) {}

  // currentDir: string = "data"
  currentDir: string = "data";
  breadcrumbList: string[] = ['My Drive'];

  dirEntries: FsEntry[] = [];
  selectedEntry: FsEntry | null = null;

  singleClick(entry: FsEntry) {
    this.selectedEntry = entry;
  }

  doubleClick(entry: FsEntry) {
    console.log(`${this.currentDir}/${entry.name}`)
    if (entry.isFolder) {
      this.currentDir = `${this.currentDir}/${entry.name}`;
      this.breadcrumbList.push(entry.name);

      this.updateDirEntries(this.currentDir);
    }
  }

  goBack(): void {
    if (this.breadcrumbList.length > 1) {
      this.breadcrumbList.pop();

      // remove last directory
      const lastIndex = this.currentDir.lastIndexOf("/");
      this.currentDir = this.currentDir.substring(0, lastIndex);

      this.updateDirEntries(this.currentDir);
    }
  }

  ngOnInit() {
    this.updateDirEntries(this.currentDir);
  }

  async updateDirEntries(path: string) {
    this.dirEntries = await this.fetchDirectoryEntries(path);
  }

  fetchDirectoryEntries(path: string): Promise<FsEntry[]> {
    const url = `http://localhost:8080/directory?path=${encodeURIComponent(path)}`;
      console.log(url)
    return this.http.get<FsEntry[]>(url).toPromise()
    .then(response => response || []);
  }
}

interface FsEntry {
  name: string;
  isFile: boolean;
  isFolder: boolean;
}
