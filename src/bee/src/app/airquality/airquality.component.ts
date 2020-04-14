import { Component, OnInit, ViewChild, AfterViewInit} from '@angular/core';
import { Title } from '@angular/platform-browser';
import {FormBuilder, FormControl, FormGroup, Validators} from '@angular/forms';
import {Air} from './air';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { catchError, map, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';


//declare var MediaRecorder: any;

@Component({
  selector: 'app-airquality',
  templateUrl: './airquality.component.html',
  styleUrls: ['./airquality.component.css']
})



export class AirqualityComponent implements OnInit {
  options: FormGroup;
  colorControl = new FormControl('primary');
  fontSizeControl = new FormControl(16, Validators.min(10));


  constructor(fb: FormBuilder, private http: HttpClient) {
    
    this.options = fb.group({
      color: this.colorControl,
      fontSize: this.fontSizeControl,
    });
  }

  ngOnInit(): void {
  }
  
  

  data: HttpRR[] = []

  getFontSize() {
    return Math.max(10, this.fontSizeControl.value);
  }

 
  getAir(): void {
    var begin = Date.now().toString()
    this.getHttpAir().subscribe(x => this.data.push({
        request_time: begin, 
        response_time: Date.now().toString(),
        duration: parseInt(Date.now().toString()) - parseInt(begin),
        response_headers: "",
        response: JSON.stringify(x),
      }))
    console.log(this.data)

  }
  getHttpAir(): Observable<object> {
    console.log("send req ...")
    //"http://a4947283c471d4a1dbae2b088549ac20-1381614647.ap-southeast-1.elb.amazonaws.com:9010/air/beijing"
    return this.http.get<object>("http://127.0.0.1:9010/air/beijing")
    .pipe(
      tap(_ => this.log('fetched air')),
      catchError(this.handleError<object>('getAir'))
    )

    

  }

  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      this.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }

  private log(message: string) {
    console.log(`AirService: ${message}`);
  }



  

}

export interface HttpRR {
  request_time: string;
  response_time: string;
  duration: number; 
  response_headers: string;
  response: string;
}

