import { Component, OnInit, ViewChild } from '@angular/core';
import { IntegrationService } from '../../services/integration.service';
import { AuthService } from "ng2-ui-auth";
import { ModalDirective } from "ng2-bootstrap";
import { Service } from "../../models/service";

@Component({
  selector: 'app-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
  providers: [IntegrationService]
})
export class IntegrateComponent implements OnInit {
  @ViewChild('childModal') public childModal: ModalDirective;
  public services: Service [] = [];


  constructor(private integrationService: IntegrationService, private auth: AuthService) { }

  ngOnInit() {
    this.services = this.integrationService.mockServices();
  }

  getRows() {
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

  slug = (providerName: string) => {
    return providerName
      .toLowerCase()
      .replace(/ /g,'-')
      .replace(/[^\w-]+/g,'');
  }

}
