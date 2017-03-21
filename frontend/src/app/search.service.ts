import { Injectable } from '@angular/core';
import { SearchConfig } from "./components/search/search.component";
import { Observable } from "rxjs";
import { Result } from "./models/result.interface";
import { JwtHttp } from "ng2-ui-auth";
import { Response } from "@angular/http";



@Injectable()
export class SearchService {
  private uri: string = 'http://localhost:8081';

  constructor(private http: JwtHttp) { }

  search(query: string): Observable<Object[]> {
    return this.http.get(`${this.uri}/search?q=${query}`)
      .map((res: Response) => res.json().data)
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  query = (query: string, config: SearchConfig): Observable<Result[]> => {
    console.log(`Searching for ${query}`);

    let gmailRS = this.mockGmailResults();
    let googleDriveRS = this.mockGoogleDriveResults();
    let slackRS = this.mockSlackResults();
    let githubRS = this.mockGithubResults();

    return new Observable(observer => {
      setTimeout(() => {
        observer.next( gmailRS);
      }, 200);

      setTimeout(() => {
        observer.next( googleDriveRS);
      }, 400);

      setTimeout(() => {
        observer.next( slackRS);
      }, 600);

      setTimeout(() => {
        observer.next( githubRS);
      }, 800);

      setTimeout(() => {
        observer.complete();
      }, 1000);
    });
  }

  private mockGmailResults(): Result[] {
    return [
      {
        account_id: "123123123",
        service: "Gmail",
        resource: "Message",
        title: "Tuesday appointment with CTO",
        description: "We should discuss tech stack for our next project",
        url: "https://mail.google.com/"
      },
      {
        account_id: "123123123",
        service: "Gmail",
        resource: "Thread",
        title: "Topditop feature discussion",
        description: "As discussed in our last meeting...",
        url: "https://mail.google.com/"
      }
    ];
  }

  private mockSlackResults(): Result[] {
    return [
      {
        account_id: "123123123",
        service: "Slack",
        resource: "Channel",
        title: "Development",
        description: "Intended for all members of FCI and CodeBehind",
        url: "https://foundcenter.slack.com/messages/development/"
      },
      {
        account_id: "123123123",
        service: "Slack",
        resource: "Mpim",
        title: "Topditop is now live",
        description: "Newest version is up, check it out",
        url: "https://foundcenter.slack.com/archives/topditop/p1489409529254206"
      },
      {
        account_id: "123123123",
        service: "Slack",
        resource: "User",
        title: "Marko Arsic",
        description: "Senior Mobile Developer",
        url: "https://foundcenter.slack.com/messages/@marsic/"
      }
    ];
  }

  private mockGithubResults(): Result[] {
    return [
      {
        account_id: "123123123",
        service: "Github",
        resource: "Commit",
        title: "TOP-152 Fix header",
        description: null,
        url: "https://github.com/foundcenter/topditop/commit/b29019f8c64d9a57dd864f72e21e77b194a75693"
      },
      {
        account_id: "123123123",
        service: "Github",
        resource: "Repository",
        title: "TopDiTop",
        description: "Newest version is up, check it out",
        url: "https://github.com/foundcenter/topditop"
      }
    ];
  }


  private mockGoogleDriveResults(): Result[] {
    return [
      {
        account_id: "54233123",
        service: "Google Drive",
        resource: "Directory",
        title: "Packator design files",
        description: "Contains all packator PDFs and UI mocks",
        url: "https://drive.google.com/drive/#my-drive"
      },
      {
        account_id: "4324222",
        service: "Google Drive",
        resource: "File",
        title: "Packator project specification",
        description: "Version 1.2 of project spec with features and sprint plans",
        url: "https://drive.google.com/drive/#my-drive"
      }
    ];
  }
}