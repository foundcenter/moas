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

  constructor(public name: string, public status: string,public  logo: string) {
  }

  getLogoUrl = () => {
    return 'assets/images/'+this.logo+'_logo.png';
  }
}

export class Account {

  constructor(public email: string) {
  }
}