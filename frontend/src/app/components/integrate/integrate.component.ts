import { Component, OnInit, ViewChild, AfterContentInit } from '@angular/core';
import { IntegrationService } from '../../services/integration.service';
import { AuthService } from 'ng2-ui-auth';
import { ModalDirective } from 'ng2-bootstrap';
import { Service } from '../../models/service';
import { AccountService } from '../../services/account.service';
import { Account } from '../../models/account';

@Component({
  selector: 'app-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
  providers: [IntegrationService, AccountService]
})
export class IntegrateComponent implements OnInit, AfterContentInit {
  @ViewChild('childModal') public childModal: ModalDirective;
  public services: Service[] = [];
  public accounts: Account[] = [];
  constructor(private integrationService: IntegrationService, private auth: AuthService, private accountService: AccountService) { }

  ngOnInit() {
    this.accounts = this.accountService.getAccounts();
    this.services = this.integrationService.getAvailableServices();
  }

  ngAfterContentInit(): void {
    this.sortAccountsToServices();
  }


  sortAccountsToServices(): void {
    this.accounts.forEach((account: Account) => {
      let service = this.findServiceByName(account.service);
      service.accounts.push(account);
    })
  }

  findServiceByName(service: Service): Service {
    let result = null;
    this.services.forEach((current: Service) => {
      if (service.name == current.name) {
        result = current;
      }
    });
    return result;
  }

  getRows(): Service[][] {
    let totalAdded = 0;
    let rows : Service[][] = [];
    let row: Service[] = [];
    for (let i = 0; i < this.services.length; i = i + 3) {
      row = [this.services[i]];
      if (this.services[i+1]) {
        row.push(this.services[i+1]);
      }
      if (this.services[i+2]) {
        row.push(this.services[i+2]);
      }
      rows.push(row);
    }

    return rows;
  }

  public showChildModal():void {
    this.childModal.show();
  }

  public hideChildModal():void {
    this.childModal.hide();
  }

  handle(serviceName: string) {
    switch (serviceName) {
      case 'gmail':
      case 'google-drive':
      case 'slack':
      case 'github':
        console.log('integrating ' + serviceName);
        this.auth.authenticate(serviceName)
          .subscribe(
            data => {
              console.log(`Data of ${serviceName} auth`);
              console.log(data);
            },
            error => {
              console.log(`Error of ${serviceName} auth`);
              console.log(error);
            },
            () => {
              console.log(`Finally of ${serviceName} auth`)
            }
          );
        break;
        
      case 'jira':
        console.log('trigger modal for jira now');
        this.showChildModal();
        break;

      default:
        console.log('Integration not handled for ' + serviceName);
    }
  }

  slug(providerName: string): string{
    return providerName
      .toLowerCase()
      .replace(/ /g,'-')
      .replace(/[^\w-]+/g,'');
  }

}
