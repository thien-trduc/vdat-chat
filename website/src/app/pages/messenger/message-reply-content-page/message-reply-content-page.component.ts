import {
  AfterViewInit,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
  ViewChild
} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {User} from '../../../core/models/user.model';
import {Message} from '../../../core/models/message/message.model';
import {CdkVirtualScrollViewport} from '@angular/cdk/scrolling';
import {MessageDataSource} from '../../../core/models/datasources/message.datasource';
import {BehaviorSubject} from 'rxjs';
import {MessageService} from '../../../core/services/ws/message.service';

@Component({
  selector: 'app-message-reply-content-page',
  templateUrl: './message-reply-content-page.component.html',
  styleUrls: ['./message-reply-content-page.component.scss']
})
export class MessageReplyContentPageComponent implements OnInit, OnChanges, AfterViewInit {

  @ViewChild('virtualScrollViewport') virtualScrollViewport: CdkVirtualScrollViewport;

  @Input() currentGroup: Group;
  @Input() currentUser: User;
  @Input() parentMessage: Message;

  @Output() collapseInfoTab = new EventEmitter();
  @Output() sortListGroupEvent = new EventEmitter();

  public messageDataSource: MessageDataSource;
  public visibleReplyMessageDrawer = new BehaviorSubject<boolean>(false);

  constructor(private messageService: MessageService) {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!this.currentGroup && !!this.currentUser && !!this.parentMessage) {
      this.messageDataSource = new MessageDataSource(this.messageService, this.currentUser,
        this.currentGroup, this.parentMessage);

      if (!!this.messageDataSource && !!this.virtualScrollViewport) {
        this.messageDataSource.virtualScroll = this.virtualScrollViewport;
      }

      if (!!this.messageDataSource && !!this.virtualScrollViewport) {
        this.messageDataSource.virtualScroll = this.virtualScrollViewport;
      }
    }
  }

  ngOnInit(): void {
  }

  ngAfterViewInit(): void {
    if (!!this.messageDataSource && !!this.virtualScrollViewport) {
      this.messageDataSource.virtualScroll = this.virtualScrollViewport;
    }
  }
}
