import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { AuthService } from "./services/auth.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit{
  title = 'MOAS app works!';

  constructor(private authService: AuthService, private router: Router) {
  }

  ngOnInit(): void {
    if (!this.authService.isLoggedIn()) {
      this.router.navigateByUrl('/login');
    } else {
      this.authService.check()
        .subscribe(
          (data) => {},
          (error) => {
            this.authService.logout();
            this.router.navigateByUrl('/login');
          }
        );
    }
  }
}
