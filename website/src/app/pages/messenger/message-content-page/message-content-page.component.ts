import {
  Component,
  Input,
  OnInit,
  Output,
  EventEmitter,
  AfterViewInit,
  OnChanges,
  SimpleChanges,
  ViewChild, ChangeDetectionStrategy
} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {MessageService} from '../../../core/services/ws/message.service';
import * as _ from 'lodash';
import {Message} from '../../../core/models/message/message.model';
import {CdkVirtualScrollViewport} from '@angular/cdk/scrolling';
import {User} from '../../../core/models/user.model';
import {BehaviorSubject} from 'rxjs';
import {MessageDataSource} from '../../../core/models/datasources/message.datasource';

@Component({
  selector: 'app-message-content-page',
  templateUrl: './message-content-page.component.html',
  styleUrls: ['./message-content-page.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MessageContentPageComponent implements OnInit, OnChanges, AfterViewInit {

  @ViewChild('virtualScrollViewport') virtualScrollViewport: CdkVirtualScrollViewport;

  @Input() currentGroup: Group;
  @Input() currentUser: User;

  @Output() collapseInfoTab = new EventEmitter();
  @Output() sortListGroupEvent = new EventEmitter();

  public messageDataSource: MessageDataSource;
  public currentMessage: Message;
  public visibleReplyMessageDrawer = new BehaviorSubject<boolean>(false);

  constructor(private messageService: MessageService) {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes && !!this.currentGroup && !!this.currentUser) {
      this.messageDataSource = new MessageDataSource(this.messageService, this.currentUser, this.currentGroup);

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

  // region Event
  public onCollapseInfoTab(): void {
    this.collapseInfoTab.emit();
  }

  public onOpenReplyMessageDrawer(message: Message): void {
    if (!!message) {
      this.currentMessage = _.cloneDeep(message);
      this.visibleReplyMessageDrawer.next(true);
    }
  }

  public onCloseReplyMessageDrawer(): void {
    this.currentMessage = null;
    this.visibleReplyMessageDrawer.next(false);
  }
  // endregion
}
