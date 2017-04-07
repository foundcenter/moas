import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { SearchService } from '../../services/search.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {

  constructor(private auth: AuthService, private searchService: SearchService, private router: Router) { }

  ngOnInit() {
  }

  searchClicked(): void {
    this.searchService.emitClick();
  }

  isActive(routeUrl: string): boolean {
    let url = this.router.url;
    return routeUrl == url;
  }

}
