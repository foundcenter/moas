import { Component } from "@angular/core";
import { ToastrConfig } from "ngx-toastr";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'MOAS app works!';

  constructor(toastrConfig: ToastrConfig) {
    toastrConfig.closeButton=true;
  }
}
