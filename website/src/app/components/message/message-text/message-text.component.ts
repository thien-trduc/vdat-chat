import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-message-text',
  templateUrl: './message-text.component.html',
  styleUrls: ['./message-text.component.scss']
})
export class MessageTextComponent implements OnInit {

  @Input() message: string;
  @Input() singleLine: boolean;

  constructor() {
  }

  ngOnInit(): void {
  }

}
