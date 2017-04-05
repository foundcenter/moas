import { Response } from '@angular/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  constructor(private router: Router, private auth: AuthService, private toastr: ToastrService) {
  }

  ngOnInit() {
  }

  loginWithGoogle() {
    this.auth.login()
      .then((res: Response) => {
        this.router.navigateByUrl('search');
        if (res.status == 201) {
          this.toastr.info('Go to Integrate page to connect other services!', 'Integrate more services')
        }
      })
      .catch(() => {
        this.toastr.error('Login with Google failed!', 'Error')
      });
  }

}
