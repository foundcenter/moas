import { Account } from './account';

export class Service {
  public accounts: Account[] = [];

  constructor(public name: string, public logo: string, public info: string) {
  }

  getLogoUrl(): string{
    return 'assets/images/'+this.logo+'_logo.png';
  }

  static JIRA(): Service {
    return new Service('Jira', 'jira', 'Search projects and issues.');
  }
  static GMAIL(): Service {
    return new Service('Gmail', 'gmail', 'Search messages, threads and attachments.');
  }
  static GOOGLEDRIVE(): Service {
    return new Service('Google Drive', 'google_drive', 'Search your drive files and directories.');
  }
  static SLACK(): Service {
    return new Service('Slack', 'slack', 'Search direct messages, channels and attachments.');
  }
  static GITHUB(): Service {
    return new Service('GitHub', 'github', 'Search commits, repositories and issues.');
  }


  }
