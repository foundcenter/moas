import { Account } from './account';

export class Service {
  public accounts: Account[] = [];

  constructor(public name: string, public logo: string) {
  }

  getLogoUrl(): string{
    return 'assets/images/'+this.logo+'_logo.png';
  }

  static JIRA(): Service {
    return new Service('Jira', 'jira');
  }
  static GMAIL(): Service {
    return new Service('Gmail', 'gmail');
  }
  static GOOGLEDRIVE(): Service {
    return new Service('Google Drive', 'google_drive');
  }
  static SLACK(): Service {
    return new Service('Slack', 'slack');
  }
  static GITHUB(): Service {
    return new Service('GitHub', 'github');
  }


  }
