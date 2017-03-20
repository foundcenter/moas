import { Service } from "./service";
export class Account {
  static readonly statusOk: string = 'Ok';
  static readonly statusExpired: string = 'Expired';

  constructor(public email: string, public id: string, public status: string, public service: Service) {
  }

  isOk(): boolean{
    return this.status == Account.statusOk;
  }

  isExpired(): boolean{
    return this.status == Account.statusExpired;
  }

}