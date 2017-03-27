import { Component, OnInit, ViewChild, AfterContentInit, OnDestroy } from '@angular/core';
import { IntegrationService } from '../../services/integration.service';
import { ModalDirective } from 'ng2-bootstrap';
import { Service } from '../../models/service';
import { AccountService } from '../../services/account.service';
import { Account } from '../../models/account';
import { AuthService } from "../../services/auth.service";
import { User } from "../../models/user";
import { Subscription } from "rxjs";

@Component({
  selector: 'app-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
  providers: [IntegrationService, AccountService]
})
export class IntegrateComponent implements OnInit, OnDestroy {
  @ViewChild('childModal') public childModal: ModalDirective;
  public services: Service[] = [];
  public accounts: Account[] = [];
  public jira: {email, password, error} = {email: '', password: '', error: null};
  private currentUserSubscription: Subscription;

  constructor(private integrationService: IntegrationService, private auth: AuthService, private accountService: AccountService) {

  }

  ngOnDestroy(): void {
    this.currentUserSubscription.unsubscribe();
  }

  ngOnInit() {
    this.services = this.integrationService.getAvailableServices();

    this.currentUserSubscription = this.auth.currentUser.subscribe((user: User) => {
      if (!user) {
        return;
      }

      this.accounts = user.accounts;
      this.assignAccountsToServices();
    });
  }

  assignAccountsToServices(): void {
    this.accounts.forEach((account: Account) => {
      this.assignAccountToService(account);
    })
  }
  assignAccountToService(account: Account): void {
    let service = this.findServiceByName(account.service);
    service.accounts.push(account);
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
    this.jira.error = null;
  }

  handle(service: Service) {
    let serviceName = this.slug(service.name);

    switch (serviceName) {
      case 'gmail':
      case 'google-drive':
      case 'slack':
      case 'github':
        console.log('integrating ' + serviceName);
        this.auth.connect(serviceName)
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

  addJira(){
    this.accountService.mockAddJira(this.jira.email, this.jira.password)
      .subscribe(
        data => {
          this.assignAccountToService(<Account>data);
          this.hideChildModal();
          this.jira.email = '';
          this.jira.password = '';
        },
        error => {
          this.jira.error = error;
        }
      );
  }

  slug(providerName: string): string{
    return providerName
      .toLowerCase()
      .replace(/ /g,'-')
      .replace(/[^\w-]+/g,'');
  }

}
