import { Injectable, OnDestroy, OnInit } from '@angular/core';
import { Response } from '@angular/http';
import { AuthService as UiAuth, JwtHttp } from "ng2-ui-auth";
import { User } from "../models/user";
import { Observable, ReplaySubject } from "rxjs";
import { Router } from "@angular/router";

@Injectable()
export class AuthService{
  private user: User = null;

  public currentUser: ReplaySubject<User> = new ReplaySubject(0);

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
        this.currentUser.next(this.user);
        return data;
      });
  }

  isLoggedIn(): boolean {
    return this.uiAuth.isAuthenticated();
  }

  logout(): Observable<void>{
    return this.uiAuth.logout();
  }

  check(): Promise<Response> {
    return this.http.get('http://localhost:8081/auth/check')
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        this.currentUser.next(this.user);
        return data;
      });
  }

  connect(serviceName: string): Observable<Response> {
    return this.uiAuth.authenticate(serviceName);
  }

  setUser(response): void {
    this.user = User.fromJson(response);
  }

}
