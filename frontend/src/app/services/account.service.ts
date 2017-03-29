import { Injectable, OnDestroy } from '@angular/core';
import { Account } from '../models/account';
import { Service } from '../models/service';
import { Observable } from "rxjs";
import { Response } from "@angular/http";
import { AuthService } from "./auth.service";
import { User } from "../models/user";
import { JwtHttp } from "ng2-ui-auth";
import { environment } from "../../environments/environment";

@Injectable()
export class AccountService {
  private uri: string = environment.apiUrl;

  private currentUser: User;

  constructor(private auth: AuthService, private http: JwtHttp) {

    auth.currentUser.subscribe((user) => {
      if (!user) {
        return;
      }
      this.currentUser = user;
    });
  }

  getAccounts(): Account[] {
    return this.currentUser.accounts;
  }

  delete(serviceName: string, id: string): Observable<Response> {
    return this.http.delete(`${this.uri}/account/${serviceName.toLowerCase()}/${id}`);
  }

  mockAddJira(email: string, password: string): Observable<Account> {
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


  mockAddOauthAccount(response: Response, service: Service): Observable<Account> {
    let account = new Account('somename@gmail.com', '77665544', Account.statusOk, service);

    return new Observable(observer => {
      setTimeout(() => {
        observer.next(account);
        observer.complete();
      }, 1000)
    });
  }

  private mockAccounts(): Account[] {
    let accounts: Account[] = [];

    let gmailAcc1 = new Account('neb.vojvodic@gmail.com', '12312312', Account.statusOk, Service.GMAIL());
    let gmailAcc2 = new Account('nv@foundcenter.com', '7654234', Account.statusOk, Service.GMAIL());
    accounts.push(gmailAcc1, gmailAcc2);

    let googleDriveAcc1 = new Account('neb.vojvodic@gmail.com', '12312312', Account.statusOk, Service.GOOGLEDRIVE());
    accounts.push(googleDriveAcc1);

    let jiraAccount = new Account('neb.vojvodic@gmail.com', '4353352', Account.statusExpired, Service.JIRA());
    accounts.push(jiraAccount);

    let githubAcc1 = new Account('neb.vojvodic@gmail.com', '7878782', Account.statusOk, Service.GITHUB());
    accounts.push(githubAcc1);

    return accounts;
  }

}
