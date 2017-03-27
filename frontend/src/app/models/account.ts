import { Service } from "./service";
export class Account {
  static readonly statusOk: string = 'Ok';
  static readonly statusExpired: string = 'Expired';

  constructor(public email: string, public id: string, public status: string,
              public service: Service) {
  }

  isOk(): boolean{
    return this.status == Account.statusOk;
  }

  isExpired(): boolean{
    return this.status == Account.statusExpired;
  }

  public static fromJson(accountResponse) {
    let email = accountResponse.email;
    let id = accountResponse.id;
    let service = Service.make(accountResponse.type);
    // TODO: check a status flag when we implement it
    let status = Account.statusOk;

    return new Account(email, id, status, service);
  }

}