import { Component, OnInit } from '@angular/core';
import {Air} from './air'

@Component({
  selector: 'app-airquality',
  templateUrl: './airquality.component.html',
  styleUrls: ['./airquality.component.css']
})
export class AirqualityComponent implements OnInit {


  air: Air = {
    indexCityVHash: '00000001',
    city: 'Beijing ',
    cityCN: '北京'
    
  }

  constructor() { }

  ngOnInit(): void {
  }

}
