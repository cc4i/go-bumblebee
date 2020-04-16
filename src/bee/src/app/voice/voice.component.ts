import { Component, OnInit, ChangeDetectorRef} from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

declare var MediaRecorder: any;

@Component({
  selector: 'app-voice',
  templateUrl: './voice.component.html',
  styleUrls: ['./voice.component.css']
})
export class VoiceComponent implements OnInit {


  mediaRecorder: any
  chunks = [];
  audioFiles = [];
  disable_record_btn =false
  disable_stop_btn = true
  
  constructor(private cd: ChangeDetectorRef, private dom: DomSanitizer) { }

  ngOnInit(): void {
    navigator.mediaDevices.getUserMedia({ audio: true, video: false }).then(stream =>  {
      this.mediaRecorder = new MediaRecorder(stream);
      this.mediaRecorder.onstop = e =>  {
        console.log('data available after MediaRecorder.stop() called.');
        var blob = new Blob(this.chunks, {type: 'audio/ogg; codecs=opus'});
        this.chunks = [];
        var audioURL = URL.createObjectURL(blob);
        // audio.src = audioURL;
        this.audioFiles.push(this.dom.bypassSecurityTrustUrl(audioURL));
        console.log(audioURL);
				console.log('recorder stopped');
				this.cd.detectChanges();
      };
      this.mediaRecorder.ondataavailable = e => {
        this.chunks.push(e.data);
      };
    },
    () => {
      alert('Error capturing audio.');
    });
  }

  retrieveVoice():void {
    navigator.mediaDevices.getUserMedia({ audio: true, video: false }).then(stream =>  {

      const context = new AudioContext();
      const source = context.createMediaStreamSource(stream);
      const processor = context.createScriptProcessor(1024, 1, 1);
      source.connect(processor);
      processor.connect(context.destination);
      
      processor.onaudioprocess = function(e) {
        // Do something with the data, e.g. convert it to WAV
        console.log(e.inputBuffer);
      };

    });
    
  }

  startRecording() {
    this.mediaRecorder.start();
    this.disable_record_btn = true
    this.disable_stop_btn = false 
		console.log(this.mediaRecorder.state);
		console.log('recorder started');
	}
	stopRecording() {
    this.mediaRecorder.stop();
    this.disable_stop_btn = true
    this.disable_record_btn = false
		console.log(this.mediaRecorder.state);
		console.log('recorder stopped');
	}

}
