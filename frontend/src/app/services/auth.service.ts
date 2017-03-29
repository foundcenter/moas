import { Injectable, OnDestroy, OnInit } from '@angular/core';
import { Response } from '@angular/http';
import { AuthService as UiAuth, JwtHttp } from "ng2-ui-auth";
import { User } from "../models/user";
import { Observable, ReplaySubject } from "rxjs";
import { Router } from "@angular/router";
import { environment } from "../../environments/environment";

@Injectable()
export class AuthService {
  private user: User = null;
  private uri: string = environment.apiUrl;

  public currentUser: ReplaySubject<User> = new ReplaySubject(1);

  constructor(private uiAuth: UiAuth, private http: JwtHttp, private router: Router) {
    if (!this.isLoggedIn()) {
      router.navigateByUrl('/login');
      return;
    }
    this.check()
      .catch(() => {
        this.logout();
        router.navigateByUrl('/login');
      })
  }

  login(): Promise<Response> {
    return this.uiAuth.authenticate('google')
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        return data;
      });
  }

  isLoggedIn(): boolean {
    return this.uiAuth.isAuthenticated();
  }

  logout(): Observable<void> {
    return this.uiAuth.logout();
  }

  check(): Promise<Response> {
    return this.http.get(`${this.uri}/auth/check`)
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        return data;
      });
  }

  connect(serviceName: string): Observable<Response> {
    return this.uiAuth.authenticate(serviceName);
  }

  setUser(response): void {
    this.user = User.fromJson(response);
    this.currentUser.next(this.user);
  }

  connectJira(url: string, username: string, password: string): Observable<Response> {
    let body = {
      'url': url,
      'username': username,
      'password': password

    }
    return this.http.post(this.uri+'/connect/jira', body)
  }

}
