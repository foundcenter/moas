import { Component, OnInit, Input } from '@angular/core';
import { Result } from "../../../models/result.interface";

@Component({
  selector: 'search-result',
  templateUrl: './result.component.html',
  styleUrls: ['./result.component.scss']
})
export class ResultComponent implements OnInit {
  @Input() result: Result;

  constructor() { }

  ngOnInit() {
  }

}
