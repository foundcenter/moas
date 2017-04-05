import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { JwtHttp } from 'ng2-ui-auth';
import { Response } from '@angular/http';
import { environment } from '../../environments/environment';


@Injectable()
export class SearchService {
  private uri: string = environment.apiUrl;

  constructor(private http: JwtHttp) { }

  search(query: string): Observable<Object[]> {
    return this.http.get(`${this.uri}/search?q=${query}`)
      .map((res: Response) => res.json().data)
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
}