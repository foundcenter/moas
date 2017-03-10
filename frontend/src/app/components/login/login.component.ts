import { Component, OnInit } from '@angular/core';
import { AuthService } from 'ng2-ui-auth';
import { Router } from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  constructor( private router: Router, private auth: AuthService ) { }

  ngOnInit() {
  }

  loginWithGoogle() {
        this.auth.authenticate('google')
            .subscribe({
                complete: () => this.router.navigateByUrl('home'),
                error: (err: any) => console.log(err)
            });
    }

}
