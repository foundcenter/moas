import { Component, OnInit, EventEmitter, AfterViewInit } from '@angular/core';

@Component({
  selector: 'a pp-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit, AfterViewInit {
  // protected providers { name: string, search: boolean}[];
  public providers : Provider[] = [];
  public configureProviders: boolean = false;
  public focusTriggeringEventEmmiter = new EventEmitter<boolean>();

  constructor() {
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
}

class Provider {
  constructor(public name: string, public search: boolean) {
  }
  toggle = () => {
    console.log(this.name + " is in provider class");
    this.search = !this.search;
  }
}
