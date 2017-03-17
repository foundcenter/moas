import { Component, OnInit } from '@angular/core';
import { IntegrationService } from "../../services/integration.service";

@Component({
  selector: 'app-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
  providers: [IntegrationService]
})
export class IntegrateComponent implements OnInit {
  public services: Service [] = [];

  constructor(private integrationService: IntegrationService) { }

  ngOnInit() {
    this.services = this.integrationService.mockServices();
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
  static readonly statusOk: string = "Ok";
  static readonly statusExpired: string = "Expired";

  constructor(public email: string, public id: string, public status: string) {
  }

  isOk = () => {
    return this.status == Account.statusOk;
  }

  isExpired = () => {
    return this.status == Account.statusExpired;
  }

}