import { Component, OnInit } from '@angular/core';
import { AuthService } from "ng2-ui-auth";
import { Router } from "@angular/router";

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.scss']
})
export class LogoutComponent implements OnInit {

  constructor(private router: Router,private authService: AuthService) { }

  ngOnInit() {
    this.authService.logout()
      .subscribe({
        complete: () => this.router.navigateByUrl('login'),
        error: (err: any) => console.log(err)
      });
  }

}
