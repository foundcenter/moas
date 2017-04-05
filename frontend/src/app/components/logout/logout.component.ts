import { ToastrService } from 'ngx-toastr';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.scss']
})
export class LogoutComponent implements OnInit {

  constructor(private router: Router, private auth: AuthService, private toastr: ToastrService) { }

  ngOnInit() {
    this.auth.logout()
      .then(() => {
        this.router.navigateByUrl('login');
        this.toastr.info('Successfully logout! See you soon', 'Success');
      })
      .catch(() => {
        this.toastr.error('Something gone wrong with logout!', 'Error');
      });
  }

}
