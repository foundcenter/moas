import { Injectable } from '@angular/core';
import { Account } from '../models/account';
import { Response } from '@angular/http';
import { AuthService } from './auth.service';
import { User } from '../models/user';
import { JwtHttp } from 'ng2-ui-auth';
import { environment } from '../../environments/environment';

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

  delete(serviceName: string, id: string): Promise<Response> {
    return this.http.delete(`${this.uri}/account/${serviceName.toLowerCase()}/${id}`)
      .toPromise()
      .then((data: Response) => {
        this.auth.setUser(data.json().data.user);
        return data;
      });
  }

}
