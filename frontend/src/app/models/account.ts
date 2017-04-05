import { Service } from "./service";
export class Account {
  static readonly statusOk: string = 'Ok';
  static readonly statusExpired: string = 'Expired';
  public email: string;
  public id: string;
  public status: string;
  public service: Service;

  constructor() {
  }

  isOk(): boolean{
    return this.status == Account.statusOk;
  }

  isExpired(): boolean{
    return this.status == Account.statusExpired;
  }

  public static fromJson(accountResponse) {
    let account = new Account();

    account.email = accountResponse.email;
    account.id = accountResponse.id;
    account.service = Service.make(accountResponse.type);
    // TODO: check a status flag when we implement it
    account.status = accountResponse.active? Account.statusOk:Account.statusExpired;

    return account;
  }

}