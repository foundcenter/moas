import { ToastrService } from 'ngx-toastr';
import { Component, OnInit } from '@angular/core';
import { AuthService } from 'ng2-ui-auth';
import { Router } from '@angular/router';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.scss']
})
export class LogoutComponent implements OnInit {

  constructor(private router: Router, private authService: AuthService, private toastrService: ToastrService) { }

  ngOnInit() {
    this.authService.logout()
      .subscribe({
        complete: () => {
          this.router.navigateByUrl('login');
          this.toastrService.info('Successfully logout! See you soon', 'Success');
        },
        error: (err: any) => {
          console.log(err);
          this.toastrService.error('Something gone wrong with logout!', 'Error');
        }
      });
  }

}
