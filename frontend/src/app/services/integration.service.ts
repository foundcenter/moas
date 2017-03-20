import { Injectable } from '@angular/core';
import { Account } from "../models/account";
import { Service } from "../models/service";

@Injectable()
export class IntegrationService {

  constructor() { }

  mockServices = () => {
    let googleAccount1 = new Account("neb.vojvodic@gmail.com", "12312312", Account.statusOk);
    let googleAccount2 = new Account("nv@foundcenter.com", "7654234", Account.statusOk);
    let gmail = new Service("Gmail", "gmail");
    gmail.accounts.push(googleAccount1, googleAccount2);

    let googleDrive = new Service("Google Drive", "google_drive");
    googleDrive.accounts.push(googleAccount1, googleAccount2);

    let atlassianAccount = new Account("neb.vojvodic@gmail.com", "4353352", Account.statusExpired);
    let jira = new Service("Jira", "jira");
    jira.accounts.push(atlassianAccount);

    let slack = new Service("Slack", "slack");

    let githubAccount1 = new Account("neb.vojvodic@gmail.com", "7878782", Account.statusOk);
    let github = new Service("Github", "github");
    github.accounts.push(githubAccount1);

    return [gmail, googleDrive, jira, slack, github];
  }
}
