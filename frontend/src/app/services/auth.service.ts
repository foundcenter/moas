import { Injectable } from '@angular/core';
import { Response } from '@angular/http';
import { AuthService as UiAuth } from "ng2-ui-auth";

@Injectable()
export class AuthService {
  public user: Object = null;

  constructor(private uiAuth: UiAuth) { }

  login(): Promise<Response> {
    return this.uiAuth.authenticate('google')
      .toPromise()
      .then((data: Response) => {
        this.user = data.json().data.user;
        return data;
      });
  }

}
