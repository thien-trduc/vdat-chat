import {Component, Input, OnInit} from '@angular/core';
import {Message} from '../../../core/models/message/message.model';

@Component({
  selector: 'app-message-file',
  templateUrl: './message-file.component.html',
  styleUrls: ['./message-file.component.scss']
})
export class MessageFileComponent implements OnInit {

  @Input() message: Message;

  constructor() {
  }

  ngOnInit(): void {
  }
}
