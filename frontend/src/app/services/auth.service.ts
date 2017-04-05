import { Injectable } from '@angular/core';
import { Response } from '@angular/http';
import { AuthService as UiAuth, JwtHttp } from 'ng2-ui-auth';
import { User } from '../models/user';
import { ReplaySubject } from 'rxjs';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment';

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
        this.logout()
          .then(() => {
            router.navigateByUrl('/login');
          });
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

  logout(): Promise<void> {
    return this.uiAuth.logout().toPromise();
  }

  check(): Promise<Response> {
    return this.http.get(`${this.uri}/auth/check`)
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        return data;
      });
  }

  connect(serviceName: string): Promise<Response> {
    return this.uiAuth.authenticate(serviceName)
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        return data;
      });
  }

  setUser(response): void {
    this.user = User.fromJson(response);
    this.currentUser.next(this.user);
  }

  setGithubPersonal(username: string, token: string): Promise<Response>  {
    return this.http.put(`${this.uri}/connect/github/${username}`, { token: token })
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        return data;
      });
  }

  connectJira(url: string, username: string, password: string): Promise<Response> {
    let body = {
      'url': url,
      'username': username,
      'password': password
    };

    return this.http.post(this.uri+'/connect/jira', body)
      .toPromise()
      .then((data: Response) => {
        this.setUser(data.json().data.user);
        return data;
      });
  }

}
