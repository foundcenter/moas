import { Component, OnInit, EventEmitter, AfterViewInit } from '@angular/core';
import { SearchService } from "../../search.service";
import { Result } from "../../interfaces/result.interface";

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
  providers: [SearchService]
})
export class SearchComponent implements OnInit, AfterViewInit {
  // protected providers { name: string, search: boolean}[];
  public providers : Provider[] = [];
  public configureProviders: boolean = false;
  public focusTriggeringEventEmmiter = new EventEmitter<boolean>();
  public results: Result[] = [];
  public presenting: string = "all";
  public resultsServices: string[] = [];
  public query: string = "";

  constructor(private searchService: SearchService) {
  }

  toggleConfigureProviders = () => {
    this.configureProviders = !this.configureProviders;
  }

  ngOnInit() {
    let gmail = new Provider("Gmail", true);
    let gdrive = new Provider("Google Drive", false);
    let slack = new Provider("Slack", true);
    let jira = new Provider("Jira", false);
    let github = new Provider("Github", true);

    this.providers.push(gmail, gdrive, slack, jira, github);

    this.autoSearch();
  }

  autoSearch() {
    this.query = "Topditop, packator";
    this.search();
  }

  ngAfterViewInit() {
    this.focusSearchBar();
  }
  
  focusSearchBar = () => {
    this.focusTriggeringEventEmmiter.emit(true);
  }

  toggle = (provider: Provider) => {
    console.log(provider.name + " is in search method");
    provider.toggle();
  }

  viewResultsBy = (service: string) => {
    if (service != null) {
      this.presenting = "all";
    }
    this.presenting = service;
  }

  sortResults = () => {
    let all = this.results;
    let services = [];
    all.forEach((result: Result) => {
      if (!services.includes(result.service)) {
        services.push((result.service));
      }
    });
    this.resultsServices = services;
  }

  resultBy(service: string): Result[] {
    if (this.presenting == "all") {
      return this.results;
    }

    return this.results.filter((result: Result) => {
      if (result.service == service) {
        return result;
      }
    });
  }

  resultCountBy(service: string): Number {
    if (service == null) {
      return this.results.length;
    }
    return this.resultBy(service).length;
  }

  search = () => {
    this.results = [];

    this.searchService.query(this.query, new SearchConfig(this.providers))
      .subscribe(
        data => {
          console.log(data);
          let result = <Result[]>data;
          this.results.push.apply(this.results, result);
          this.sortResults();
        }
      );

  }
}

export class Provider {
  constructor(public name: string, public search: boolean) {
  }
  toggle = () => {
    console.log(this.name + " is in provider class");
    this.search = !this.search;
  }
}

export class SearchConfig {
  constructor(public providers: Provider[]) {
  }
}
