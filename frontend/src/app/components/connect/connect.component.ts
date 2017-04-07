import { AfterViewInit, Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { ModalDirective } from 'ngx-bootstrap';
import { Service } from '../../models/service';
import { Response } from '@angular/http';
import { Account } from '../../models/account';
import { AuthService } from '../../services/auth.service';
import { User } from '../../models/user';
import { Subscription } from 'rxjs';
import { ConnectService } from '../../services/connect.service';

@Component({
  selector: 'app-connect',
  templateUrl: './connect.component.html',
  styleUrls: ['./connect.component.scss'],
  providers: [ConnectService]
})
export class ConnectComponent implements OnInit, OnDestroy {
  @ViewChild('childModal') public childModal: ModalDirective;
  public focusTriggeringEventEmmiter = new EventEmitter<boolean>();
  public services: Service[] = [];
  public accounts: Account[] = [];
  public rows: Service[][] = [];
  public jira: { username: string, password: string, url: string, error: string, custom: boolean };
  public github: { token: string, open: boolean, error: string };
  public githubPersonalToken: string = '';
  private currentUserSubscription: Subscription;

  constructor(private connectService: ConnectService, private auth: AuthService, private toastr: ToastrService) {
    this.github = { token: '', open: false, error: '' };
    this.jira = { username: '', password: '', url: '', error: null, custom: false };
  }

  ngOnDestroy(): void {
    this.currentUserSubscription.unsubscribe();
  }

  ngOnInit() {
    this.currentUserSubscription = this.auth.currentUser.subscribe((user: User) => {
      if (!user) {
        return;
      }
      this.services = this.connectService.getAll();
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
    setTimeout(() => {
      this.focusTriggeringEventEmmiter.emit(true);
    }, 500);
  }

  public hideChildModal(): void {
    this.childModal.hide();
    this.jira.error = null;
  }

  deleteAccount(account: Account, service: Service) {
    this.auth.delete(service.name, account.id)
      .then(() => {
        this.toastr.success(`Successfully deleted ${service.name} account ${account.id}!`, 'Account deleted.');
      })
      .catch(() => {
        this.toastr.error(`While deleting ${service.name} account ${account.id}!`, 'Account not deleted.');
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
        this.github.token = '';
        this.toastr.success(`Personal token for ${githubUsername} has been saved.`, 'Github updated.');
      })
      .catch((error: Response) => {
        this.toastr.error(`${error.json().error}`, `Personal token error.`);
        this.github.token = '';
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
            this.toastr.success(`${serviceName} account successfully connected.`, 'Account added.')
          })
          .catch(() => {
            this.toastr.error(`${serviceName} account failed to connect.`, 'Failed to add account.')
          });
        break;

      case 'jira':
        this.showChildModal();
        break;

      default:
        console.log('Connection not handled for ' + serviceName);
    }
  }

  toggleJiraUrl(): void {
    this.jira.custom = !this.jira.custom;
  }

  private getJiraUrl(): string {
    if (this.jira.custom) {
      return this.jira.url;
    }
    return `https://${this.jira.url}.atlassian.net`;
  }
  addJira() {
    this.auth.connectJira(this.getJiraUrl(), this.jira.username, this.jira.password)
      .then(() => {
        this.jira.username = '';
        this.jira.password = '';
        this.jira.url = '';
        this.toastr.success('Jira successfully connected.', 'Account added.');
        this.hideChildModal();
      })
      .catch(error => {
        this.jira.error = error;
        this.toastr.error('Jira failed to connect.', 'Failed to add account.');
      });
  }
}
