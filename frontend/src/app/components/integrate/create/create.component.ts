import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from "@angular/router";

@Component({
  selector: 'app-create',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.scss']
})
export class CreateComponent implements OnInit {
  private service: string;

  constructor(private route: ActivatedRoute) { }

  ngOnInit() {
    this.service = this.route.snapshot.params['service'];
    console.log(this.service);
  }

}
