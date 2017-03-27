import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { AuthService } from "../../services/auth.service";

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
