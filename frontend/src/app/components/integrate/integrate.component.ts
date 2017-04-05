import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { IntegrationService } from '../../services/integration.service';
import { ModalDirective } from 'ng2-bootstrap';
import { Service } from '../../models/service';
import { AccountService } from '../../services/account.service';
import { Response } from '@angular/http';
import { Account } from '../../models/account';
import { AuthService } from '../../services/auth.service';
import { User } from '../../models/user';
import { Subscription } from 'rxjs';

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
  public rows: Service[][] = [];
  public jira: { username, password, url, error } = { username: '', password: '', url: '', error: null };
  public github: { token: string, open: boolean, error: string } = { token: '', open: false, error: '' };
  public githubPersonalToken: string = '';
  private currentUserSubscription: Subscription;

  constructor(private integrationService: IntegrationService, private auth: AuthService, private accountService: AccountService, private toastrService: ToastrService) { }

  ngOnDestroy(): void {
    this.currentUserSubscription.unsubscribe();
  }

  ngOnInit() {
    this.currentUserSubscription = this.auth.currentUser.subscribe((user: User) => {
      if (!user) {
        return;
      }
      this.services = this.integrationService.getAll();
      this.accounts = user.accounts;
      this.rows = this.getRows();
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
    let rows: Service[][] = [];
    let row: Service[] = [];

    for (let i = 0; i < this.services.length; i = i + 3) {
      row = [this.services[i]];
      if (this.services[i + 1]) {
        row.push(this.services[i + 1]);
      }
      if (this.services[i + 2]) {
        row.push(this.services[i + 2]);
      }
      rows.push(row);
    }

    return rows;
  }

  public showChildModal(): void {
    this.childModal.show();
  }

  public hideChildModal(): void {
    this.childModal.hide();
    this.jira.error = null;
  }

  deleteAccount(account: Account, service: Service) {
    this.accountService.delete(service.name, account.id)
      .then(() => {
        this.toastrService.success(`Successfully delete ${service.name} account ${account.id} !`, 'Account deleted');
      })
      .catch(() => {
        this.toastrService.error(`while deleting ${service.name} account ${account.id} !`, 'Error');
      });
  }

  openGithub(account: Account): void{
    this.github.open = true;
  }
  closeGithub(): void {
    this.github.open = false;
    this.github.token = '';
    this.github.error = null;
  }
  submitPersonal(githubUsername: string) {
    this.auth.setGithubPersonal(githubUsername, this.github.token)
      .then(() => {
        this.github.open = false;
        this.github.token = "";
        this.toastrService.success(`Personal token for ${githubUsername} has been saved.`, 'Github updated');
      })
      .catch((error: Response) => {
        this.toastrService.error(`${error.json().error}`, `Personal token error`);
        this.github.token = "";
      });
  }

  handle(service: Service) {
    let serviceName = service.slug();

    switch (serviceName) {
      case 'gmail':
      case 'drive':
      case 'slack':
      case 'github':
        this.auth.connect(serviceName)
          .then(() => {
            this.toastrService.success(`${serviceName} successfuly connected`, 'Success')
          })
          .catch(() => {
            this.toastrService.error(`${serviceName} fail to connect`, 'Error')
          });
        break;

      case 'jira':
        this.showChildModal();
        break;

      default:
        console.log('Integration not handled for ' + serviceName);
    }
  }

  addJira() {
    this.auth.connectJira(this.jira.url, this.jira.username, this.jira.password)
      .then(() => {
        this.jira.username = '';
        this.jira.password = '';
        this.jira.url = '';
        this.toastrService.success('Jira successfuly connected', 'Success');
        this.hideChildModal();
      })
      .catch(error => {
        this.jira.error = error;
        this.toastrService.error('Jira fail to connect', 'Error');
      });
  }
}
