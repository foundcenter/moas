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

  search = () => {
    this.results = [];

    this.searchService.query(this.query, new SearchConfig(this.providers))
      .subscribe(
        data => {
          console.log(data);
          let result = <Result[]>data;
          this.results.push.apply(this.results, result);
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
