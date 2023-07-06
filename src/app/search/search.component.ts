import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.css']
})
export class SearchComponent {
  constructor(private http: HttpClient) {}

  searchTerm: string = "";
  searchResult: searchRes | null = null;

  search() {
    const index = "hpotter.bleve";
    const url = `http://localhost:8080/search?i=${encodeURIComponent(index)}&q=${encodeURIComponent(this.searchTerm)}`;
    console.log(url)
    this.http.get<searchRes>(url).subscribe(
      (result: searchRes) => {
        console.log(result);
        this.searchResult = result
      },
      (error: any) => {
        console.error('Error:', error);
      }
    );
  }
}

interface SearchHit {
  Name: string;
  Line: string[];
}

interface searchRes {
  SearchStat: string;
  Hits:       SearchHit[];
}
