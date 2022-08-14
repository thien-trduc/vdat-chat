import {Message} from '../message/message.model';
import {CollectionViewer} from '@angular/cdk/collections';
import {MessageService} from '../../services/ws/message.service';
import {filter, finalize, map, takeUntil} from 'rxjs/operators';
import * as _ from 'lodash';
import {CdkVirtualScrollViewport} from '@angular/cdk/scrolling';
import {Group} from '../group/group.model';
import {User} from '../user.model';
import {BaseDataSource} from './base.datasource';

export class MessageDataSource extends BaseDataSource<Message> {
  public MAX_ITEM_PER_PAGE = 7;

  private virtualScrollViewport: CdkVirtualScrollViewport;
  private scrollTopOld: number;
  private firstLoadMessage = false;
  private lastMessage: Message;

  private audioNotification: HTMLAudioElement;

  constructor(private messageService: MessageService,
              private currentUser: User,
              private currentGroup: Group,
              private parentMessage?: Message) {
    super();
    this.pageSize = 20;
    this.messagesListener();

    this.audioNotification = new Audio('/assets/audios/pristine.mp3');
    this.audioNotification.load();
  }

  public set virtualScroll(virtualScrollViewport: CdkVirtualScrollViewport) {
    this.virtualScrollViewport = virtualScrollViewport;
    this.scrollTopOld = 100;

    this.virtualScrollViewport.elementScrolled()
      .pipe(filter(() => this.virtualScrollViewport.measureScrollOffset('top') === 0))
      .subscribe(() => {
        this.lastMessage = this.cachedData[0];
        if (!!this.lastMessage && !!this.parentMessage) {
          this.loadOldMessage(this.lastMessage.id, this.parentMessage.id);
        } else if (!!this.lastMessage) {
          this.loadOldMessage(this.lastMessage.id);
        }
      });

    this.scrollToBottom();
  }

  protected setup(collectionViewer: CollectionViewer): void {
    if (!!this.parentMessage) {
      this.loadOldMessage(null, this.parentMessage.id);
    } else {
      this.loadOldMessage();
    }

    collectionViewer.viewChange
      .pipe(
        takeUntil(this.complete$),
        takeUntil(this.disconnect$))
      .subscribe(range => {
        // if (range.start === 0) {
        //   const lastMessage: Message = this.cachedData[0];
        //   if (!!lastMessage) {
        //     this.loadOldMessage(lastMessage.id);
        //   }
        // }
      });
  }

  private loadOldMessage(lastMessageId?: number, parentMessageId?: number): void {
    if (!!this.currentGroup) {
      if (!!parentMessageId) {
        this.messageService.getReplyMessagesHistory(this.currentGroup.id, parentMessageId, lastMessageId);
      } else {
        this.messageService.getMessagesHistory(this.currentGroup.id, lastMessageId);
      }
    }
  }

  private messagesListener(): void {
    this.messageService.messages$
      .pipe(filter(messageResponse => messageResponse.groupId === this.currentGroup.id))
      .pipe(map(newMessageResponse => Message.fromJson(newMessageResponse.message)))
      .pipe(filter(message => !!this.parentMessage
        ? message.parentId === this.parentMessage.id : message.parentId === -1
      ))
      .subscribe(message => {
        this.firstLoadMessage = this.cachedData.length > 0;
        const foundMessage: Message = this.cachedData.find(iter => iter.id === message.id);

        if (!foundMessage) {
          if (!!this.currentGroup && !!this.currentGroup.members) {
            const sender = message.sender;
            message.isOwner = sender.userId === this.currentUser.userId;
          }

          this.cachedData.push(message);
          this.cachedData = _.sortBy(this.cachedData, 'id');

          this.dataStream.next(this.cachedData);

          if (!!this.lastMessage) {
            this.scrollToMessage(this.lastMessage.id);
          }
        }
      });

    this.messageService.newMessage$
      .pipe(filter(newMessageResponse => newMessageResponse.groupId === this.currentGroup.id))
      .pipe(map(newMessageResponse => {
        const message = Message.fromJson(newMessageResponse.message);

        if (!!this.currentGroup && !!this.currentGroup.members) {
          const sender = message.sender;
          message.isOwner = sender.userId === this.currentUser.userId;
        }

        // cập nhật số lượng tin nhắn con
        if (message.parentId !== -1) {
          const foundParentMessage = this.cachedData.find(iter => iter.id === message.parentId);
          if (!!foundParentMessage) {
            foundParentMessage.totalChildMessage += 1;

            console.log(message);

            if (foundParentMessage.isOwner && !message.isOwner) {
              this.notifyNewMessage();
            }
          }
        }

        return message;
      }))
      .pipe(filter(message => !!this.parentMessage
        ? message.parentId === this.parentMessage.id : message.parentId === -1
      ))
      .subscribe(message => {
        const foundMessage: Message = this.cachedData.find(iter => iter.id === message.id);

        if (!foundMessage) {
          this.cachedData.splice(
            0,
            this.cachedData.length <= this.MAX_ITEM_PER_PAGE ? 0 : 1,
            message);
          this.cachedData = _.sortBy(this.cachedData, 'id');

          this.dataStream.next(this.cachedData);

          if (!message.isOwner) {
            this.notifyNewMessage();
          } else {
            this.scrollToBottom();
          }
        }
      });

    this.messageService.updateMessage$
      .pipe(filter(updateMessageResponse => updateMessageResponse.groupId === this.currentGroup.id))
      .pipe(map(updateMessageResponse => Message.fromJson(updateMessageResponse.message)))
      .pipe(filter(message => !!this.parentMessage
        ? message.parentId === this.parentMessage.id : message.parentId === -1
      ))
      .subscribe(message => console.log(message));

    this.messageService.deleteMessage$
      .pipe(filter(deleteMessageResponse => deleteMessageResponse.groupId === this.currentGroup.id))
      .pipe(map(deleteMessageResponse => deleteMessageResponse.messageId))
      .subscribe(messageId => {
        const message = this.cachedData.find(iter => iter.id === messageId);
        message.deletedAt = new Date();
        message.message = null;
        this.dataStream.next(this.cachedData);
      });
  }

  private notifyNewMessage(): void {
    if (!this.audioNotification.ended) {
      this.audioNotification.pause();
      this.audioNotification.currentTime = 0;
    }
    this.audioNotification.play().then();
  }

  public scrollToBottom(): void {
    if (this.virtualScrollViewport) {
      setTimeout(() => {
        requestAnimationFrame(() => {
          this.virtualScrollViewport.scrollTo({
            behavior: 'smooth',
            bottom: 0
          });
        });
      }, 250);
    }
  }

  private scrollToMessage(messageId: number): void {
    if (this.virtualScrollViewport) {
      const index = this.cachedData.findIndex(iter => iter.id === messageId);

      console.log(index);

      requestAnimationFrame(() => {
        this.virtualScrollViewport.scrollToIndex(index, 'auto');
      });
    }
  }
}
