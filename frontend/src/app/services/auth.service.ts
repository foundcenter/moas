import { Injectable } from '@angular/core';
import { Response } from '@angular/http';
import { AuthService as UiAuth, JwtHttp } from "ng2-ui-auth";
import { User } from "../models/user";
import { Observable } from "rxjs";

@Injectable()
export class AuthService {
  private user: User = null;

  constructor(private uiAuth: UiAuth, private http: JwtHttp) { }

  login(): Promise<Response> {
    return this.uiAuth.authenticate('google')
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        console.log(this.user);
        return data;
      });
  }

  isLoggedIn(): boolean {
    return this.uiAuth.isAuthenticated();
  }

  logout(): Observable<void>{
    return this.uiAuth.logout();
  }

  check(): Observable<Response> {
    return this.http.get('http://localhost:8081/auth/check');
  }

  connect(serviceName: string): Observable<Response> {
    return this.uiAuth.authenticate(serviceName);
  }

  setUser(response): void {
    this.user = User.fromJson(response);
  }

  public getUser(): User {
    return this.user;
  }

}
