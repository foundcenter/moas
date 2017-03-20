import { Injectable } from '@angular/core';
import { Service } from '../models/service';

@Injectable()
export class IntegrationService {

  constructor() { }

  getAvailableServices(): Service[] {
    return this.mockServices();
  }

  private mockServices(): Service[] {
    let gmail = new Service('Gmail', 'gmail');
    let googleDrive = new Service('Google Drive', 'google_drive');
    let jira = new Service('Jira', 'jira');
    let slack = new Service('Slack', 'slack');
    let github = new Service('Github', 'github');

    return [gmail, googleDrive, jira, slack, github];
  }
}
