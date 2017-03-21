import { Component, OnInit, EventEmitter, AfterViewInit } from '@angular/core';
import { SearchService } from '../../search.service';
import { Result } from "../../models/result.interface";

const ALL = 'All';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
  providers: [SearchService]
})
export class SearchComponent implements OnInit, AfterViewInit {
  public focusTriggeringEventEmmiter = new EventEmitter<boolean>();
  public results: Result[] = [];
  public resultsServices: string[] = [ALL];
  public selected: string = ALL;
  public query: string = '';

  constructor(private searchService: SearchService) {
  }

  ngOnInit() {

  }

  ngAfterViewInit() {
    this.focusSearchBar();
  }
  
  private focusSearchBar(): void{
    this.focusTriggeringEventEmmiter.emit(true);
  }

  private sortResults(): void{
    this.results.forEach((result: Result) => {
      if (!this.resultsServices.includes(result.service)) {
        this.resultsServices.push((result.service));
      }
    });
  }

  select(service: string): void {
    this.selected = service;
  }

  filterBy(service: string): Result[] {
    if (service == ALL) {
      return this.results;
    }
    return this.results.filter((result: Result) => {
      if (result.service == service) {
        return result;
      }
    });
  }

  countBy(service: string): Number {
    if (service == ALL) {
      return this.results.length;
    }
    return this.filterBy(service).length;
  }

  search(): void{
    this.reset();

    this.searchService.search(this.query)
      .subscribe(
        data => {
          let result = <Result[]>data;
          this.results.push.apply(this.results, result);
          this.sortResults();
        }
      );
  }

  reset(): void {
    this.results = [];
    this.resultsServices = [ALL];
    this.selected = ALL;
  }
}

