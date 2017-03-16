import { Injectable } from '@angular/core';
import { Service } from "../components/integrate/integrate.component";

@Injectable()
export class IntegrationService {

  constructor() { }

  mockServices = () => {
    let gmail = new Service("Gmail", "Ok", "gmail");
    let googleDrive = new Service("Google Drive", "Expired", "google_drive");
    let jira = new Service("Jira", "Ok", "jira");
    let slack = new Service("Slack", "Not integrated", "slack");

    return [gmail, googleDrive, jira, slack];
  }
}
