import { Injectable } from '@angular/core';
import { JwtHttp } from 'ng2-ui-auth';
import { Response } from '@angular/http';
import { environment } from '../../environments/environment';
import { Subject } from 'rxjs';

@Injectable()
export class SearchService {
  private uri: string = environment.apiUrl;
  public searchClicked: Subject<boolean> = new Subject<boolean>();

  constructor(private http: JwtHttp) { }

  search(query: string): Promise<Response> {
    return this.http.get(`${this.uri}/search?q=${query}`)
      .toPromise();
  }

  emitClick(): void {
    this.searchClicked.next(true);
  }
}