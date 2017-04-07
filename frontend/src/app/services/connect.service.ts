import { Injectable } from '@angular/core';
import { Service } from '../models/service';

@Injectable()
export class ConnectService {

  constructor() { }

  getAll(): Service[] {
    let gmail = Service.GMAIL();
    let googleDrive = Service.GOOGLEDRIVE();
    let jira = Service.JIRA();
    let slack = Service.SLACK();
    let github = Service.GITHUB();

    return [gmail, googleDrive, jira, slack, github];
  }
}
