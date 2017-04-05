import { Injectable } from '@angular/core';
import { JwtHttp } from 'ng2-ui-auth';
import { Response } from '@angular/http';
import { environment } from '../../environments/environment';

@Injectable()
export class SearchService {
  private uri: string = environment.apiUrl;

  constructor(private http: JwtHttp) { }

  search(query: string): Promise<Response> {
    return this.http.get(`${this.uri}/search?q=${query}`)
      .toPromise();
  }
}