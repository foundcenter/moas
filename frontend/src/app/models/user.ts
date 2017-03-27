import { Account } from "./account";

export class User {
  public id: string;
  public name: string;
  public picture: string;
  public emails: string[];
  public accounts: Account[];

  public static fromJson(response){
    let user = new User();
    user.id = response.id;
    user.name = response.id;
    user.picture = response.picture;
    user.emails = response.emails;
    user.accounts = response.accounts.map((account) => {
      return Account.fromJson(account);
    });

    return user;
  }
}