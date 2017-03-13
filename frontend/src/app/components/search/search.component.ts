import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'a pp-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit {
  // protected providers { name: string, search: boolean}[];
  public providers : Provider[] = [];
  public configureProviders: boolean = false;

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

    this.providers.push(gmail, gdrive, slack, jira);
  }

  toggle = (provider: Provider) => {
    provider.toggle();
  }
}

class Provider {
  constructor(public name: string, public search: boolean) {
  }
  toggle = () => {
    this.search = !this.search;
  }
}
