import { Component, OnInit } from '@angular/core';
import { IntegrationService } from '../../services/integration.service';
import { AuthService } from "ng2-ui-auth";

@Component({
  selector: 'app-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
  providers: [IntegrationService]
})
export class IntegrateComponent implements OnInit {
  public services: Service [] = [];

  constructor(private integrationService: IntegrationService, private auth: AuthService) { }

  ngOnInit() {
    this.services = this.integrationService.mockServices();
  }

  handle = (serviceName: string) => {
    switch (serviceName) {
      case 'gmail':
      case 'google-drive':
      case 'slack':
      case 'github':
        console.log('integrating ' + serviceName);
        this.auth.authenticate(serviceName)
          .subscribe(
            data => console.log(data),
            error => console.log(error),
            () => console.log(`Integrated ${serviceName} successfuly`)
          );
        break;
        
      case 'jira':
        console.log('trigger modal for jira now');
        
        break;

      default:
        console.log('Integration not handled for ' + serviceName);
    }
  }

  slug = (providerName: string) => {
    return providerName
      .toLowerCase()
      .replace(/ /g,'-')
      .replace(/[^\w-]+/g,'');
  }

}

export class Service {
  public accounts: Account[] = [];

  constructor(public name: string, public logo: string) {
  }

  getLogoUrl = () => {
    return 'assets/images/'+this.logo+'_logo.png';
  }
}

export class Account {
  static readonly statusOk: string = 'Ok';
  static readonly statusExpired: string = 'Expired';

  constructor(public email: string, public id: string, public status: string) {
  }

  isOk = () => {
    return this.status == Account.statusOk;
  }

  isExpired = () => {
    return this.status == Account.statusExpired;
  }

}