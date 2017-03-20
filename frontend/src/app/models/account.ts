export class Account {
  static readonly statusOk: string = 'Ok';
  static readonly statusExpired: string = 'Expired';

  constructor(public email: string, public id: string, public status: string) {
  }

  isOk = () => {
    return this.status == Account.statusOk;
  }

  isExpired = () => {
    return this.status == Account.statusExpired;
  }

}