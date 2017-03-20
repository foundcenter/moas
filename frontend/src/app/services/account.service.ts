import { Injectable } from '@angular/core';
import { Account } from '../models/account';
import { Service } from '../models/service';
import { Observable } from "rxjs";

@Injectable()
export class AccountService {

  constructor() { }

  getAccounts(): Account[] {
    return this.mockAccounts();
  }

  addJira(email: string, password: string): Observable<Account> {
    if (!password) {
      return Observable.throw(new Error('Bad credentials'));
    }

    if (email == 'already.connected') {
      return Observable.throw(new Error('Already connected'));
    }

    let newJira = new Account(email, "989898", Account.statusOk, Service.JIRA());

    return new Observable((observer) => {
      setTimeout(() => {
        observer.next(newJira);
        observer.complete();
      }, 500);
    })
  }

  private mockAccounts(): Account[] {
    let accounts: Account[] = [];

    let gmailService = new Service('Gmail', 'gmail');
    let gmailAcc1 = new Account('neb.vojvodic@gmail.com', '12312312', Account.statusOk, gmailService);
    let gmailAcc2 = new Account('nv@foundcenter.com', '7654234', Account.statusOk, gmailService);
    accounts.push(gmailAcc1, gmailAcc2);

    let googleDriveService = new Service('Google Drive', 'google_drive');
    let googleDriveAcc1 = new Account('neb.vojvodic@gmail.com', '12312312', Account.statusOk, googleDriveService);
    accounts.push(googleDriveAcc1);

    let jiraService = new Service('Jira', 'jira');
    let jiraAccount = new Account('neb.vojvodic@gmail.com', '4353352', Account.statusExpired, jiraService);
    accounts.push(jiraAccount);

    let githubService = new Service('Github', 'github');
    let githubAcc1 = new Account('neb.vojvodic@gmail.com', '7878782', Account.statusOk, githubService);
    accounts.push(githubAcc1);

    return accounts;
  }

}
