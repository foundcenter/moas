import { Injectable } from '@angular/core';
import { SearchConfig, Provider } from "./components/search/search.component";
import { Observable } from "rxjs";

export interface Result {
  accountId: string;
  service: string;
  resource: string;
  title: string;
  description: string;
  url: string;
}

@Injectable()
export class SearchService {

  constructor() { }

  query = (query: string, config: SearchConfig): Observable<Result[]> => {
    console.log(`Searching for ${query}`);

    let gmailRS = this.mockGmailResults();
    let googleDriveRS = this.mockGoogleDriveResults();
    let slackRS = this.mockSlackResults();
    let githubRS = this.mockGithubResults();

    return new Observable(observer => {
      setTimeout(() => {
        observer.next( gmailRS);
      }, 500);

      setTimeout(() => {
        observer.next( googleDriveRS);
      }, 1000);

      setTimeout(() => {
        observer.next( slackRS);
      }, 1500);

      setTimeout(() => {
        observer.next( githubRS);
      }, 2000);

      setTimeout(() => {
        observer.complete();
      }, 2500);
    });
  }

  private mockGmailResults(): Result[] {
    return [
      {
        accountId: "123123123",
        service: "Gmail",
        resource: "Message",
        title: "Tuesday appointment with CTO",
        description: "We should discuss tech stack for our next project",
        url: "https://mail.google.com/"
      },
      {
        accountId: "123123123",
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
        accountId: "123123123",
        service: "Slack",
        resource: "Channel",
        title: "Development",
        description: "Intended for all members of FCI and CodeBehind",
        url: "https://foundcenter.slack.com/messages/development/"
      },
      {
        accountId: "123123123",
        service: "Slack",
        resource: "Mpim",
        title: "Topditop is now live",
        description: "Newest version is up, check it out",
        url: "https://foundcenter.slack.com/archives/topditop/p1489409529254206"
      },
      {
        accountId: "123123123",
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
        accountId: "123123123",
        service: "Github",
        resource: "Commit",
        title: "TOP-152 Fix header",
        description: null,
        url: "https://github.com/foundcenter/topditop/commit/b29019f8c64d9a57dd864f72e21e77b194a75693"
      },
      {
        accountId: "123123123",
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
        accountId: "54233123",
        service: "Google Drive",
        resource: "Directory",
        title: "Packator design files",
        description: "Contains all packator PDFs and UI mocks",
        url: "https://drive.google.com/drive/#my-drive"
      },
      {
        accountId: "4324222",
        service: "Google Drive",
        resource: "File",
        title: "Packator project specification",
        description: "Version 1.2 of project spec with features and sprint plans",
        url: "https://drive.google.com/drive/#my-drive"
      }
    ];
  }
}