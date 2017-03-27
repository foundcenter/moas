import { Component, OnInit } from '@angular/core';
import { AuthService, JwtHttp } from "ng2-ui-auth";
import { Router } from "@angular/router";
import { ToastrConfig } from 'ngx-toastr';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'MOAS app works!';

  constructor(private authService: AuthService, private router: Router, private http: JwtHttp, toastrConfig: ToastrConfig) {
    toastrConfig.closeButton=true;
  }

  ngOnInit(): void {
    if (!this.authService.isAuthenticated()) {
      this.router.navigateByUrl('/login');
    } else {
      this.http.get('http://localhost:8081/auth/check')
        .subscribe(
        (data) => { },
        (error) => {
          this.authService.logout();
          this.router.navigateByUrl('/login');
        }
        );
    }
  }
}
