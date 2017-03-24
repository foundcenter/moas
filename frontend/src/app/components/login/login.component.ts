import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { AuthService } from "../../services/auth.service";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  constructor(private router: Router, private auth: AuthService) {
  }

  ngOnInit() {
  }

  loginWithGoogle() {
    this.auth.login()
      .then(() => {
        this.router.navigateByUrl('search');
      })
      .catch(() => {
        console.log("show toastr that login failed");
      });
  }

}
