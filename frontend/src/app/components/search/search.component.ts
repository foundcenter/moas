import { AfterViewInit, Component, EventEmitter, OnInit } from '@angular/core';
import { SearchService } from '../../services/search.service';
import { Result } from '../../models/result.interface';
import { Response } from '@angular/http';
import { AccountError } from '../../models/accountError';
import { ToastrService } from 'ngx-toastr';

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
  public noResults: boolean = false;
  public errors: AccountError[] = [];
  public loading: boolean =false;

  constructor(private searchService: SearchService, private toastr: ToastrService) {
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
    if (this.results.length == 0) {
      this.noResults = true;
      return;
    }
    this.results.forEach((result: Result) => {
      if (!this.resultsServices.includes(result.service)) {
        this.resultsServices.push((result.service));
      }
    });
    this.noResults = false;
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

  onKey(event: KeyboardEvent): void {
    if (this.query.length == 0 && this.noResults) {
      this.reset();
    }
  }

  private isValidQuery(): boolean {
    return this.query.length > 0;
  }

  search(): void{
    this.reset();

    if (!this.isValidQuery()) {
      return;
    }

    this.loading=true;

    this.searchService.search(this.query)
      .then((data: Response) => {
        this.results = <Result[]>data.json().data;
        this.sortResults();

        if (<AccountError[]>data.json().meta) {
          this.errors = <AccountError[]>data.json().meta.errors;
        }
        //set flags for accounts on user object in auth service
        this.loading=false;
      })
      .catch((error) => {
        this.toastr.error('An error occurred', 'Search failed');
        this.loading=false;
      });
  }

  reset(): void {
    this.results = [];
    this.resultsServices = [ALL];
    this.selected = ALL;
    this.noResults = false;
  }
}

