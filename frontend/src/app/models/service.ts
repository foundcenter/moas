import { Account } from './account';

export class Service {
  public accounts: Account[] = [];

  constructor(public name: string, public logo: string, public info: string) {
  }

  getLogoUrl(): string {
    return 'assets/images/' + this.logo + '_logo.png';
  }

  slug(): string {
    switch (this.name) {
      case Service.GOOGLEDRIVE().name:
        return 'drive';
      default:
        return this.name.toLowerCase();
    }
  }

  static make(provider: string): Service {
    switch (provider) {
      case 'gmail':
        return Service.GMAIL();
      case 'drive':
        return Service.GOOGLEDRIVE();
      case 'github':
        return Service.GITHUB();
      case 'slack':
        return Service.SLACK();
      case 'jira':
        return Service.JIRA();
      default:
        console.log('instantiating default provider');
        return Service.GOOGLEDRIVE();
    }
  }

  static JIRA(): Service {
    return new Service('Jira', 'jira', 'Search projects and issues.');
  }

  static GMAIL(): Service {
    return new Service('Gmail', 'gmail', 'Search messages, threads and attachments.');
  }

  static GOOGLEDRIVE(): Service {
    return new Service('Google Drive', 'drive', 'Search your drive files and directories.');
  }

  static SLACK(): Service {
    return new Service('Slack', 'slack', 'Search direct messages, channels and attachments.');
  }

  static GITHUB(): Service {
    return new Service('GitHub', 'github', 'Search commits, repositories and issues.');
  }


}
