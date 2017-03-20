import { Account } from './account';

export class Service {
  public accounts: Account[] = [];

  constructor(public name: string, public logo: string) {
  }

  getLogoUrl(): string{
    return 'assets/images/'+this.logo+'_logo.png';
  }
}
