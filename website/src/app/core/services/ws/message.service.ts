import {Injectable} from '@angular/core';
import {WebSocketSubject} from 'rxjs/internal-compatibility';
import {EMPTY, from, Observable, ReplaySubject, timer} from 'rxjs';
import {DeleteMessageResponseBody, NewMessageResponseBody} from '../../models/message/message-response-body';
import {MessageResponse} from '../../models/message/message-response';
import {catchError, delay, delayWhen, filter, map, retry, retryWhen, switchAll, tap} from 'rxjs/operators';
import {MessageResponseType} from '../../models/message/message-response-type.enum';
import {webSocket} from 'rxjs/webSocket';
import {environment} from '../../../../environments/environment';
import {StorageService} from '../commons/storage.service';
import {User} from '../../models/user.model';
import {
  DeleteMessageBody,
  LoadOldMessageBody,
  SendMessageBody,
  SubscribeMessageBody,
  UpdateMessageBody
} from '../../models/message/message-request-body';
import * as _ from 'lodash';
import {MessageRequest} from '../../models/message/message-request';
import {MessageRequestType} from '../../models/message/message-request-type.enum';
import {MessageType} from '../../models/message/message-type.enum';

export const WS_ENDPOINT = `${environment.api.messages.protocol}://${window.location.host}${environment.api.messages.path}`;

@Injectable({
  providedIn: 'root'
})
export class MessageService {
  private socket$: WebSocketSubject<MessageResponse<any>>;
  private messageSubject$ = new ReplaySubject<Observable<MessageResponse<any>>>();
  private currentUser: User;

  public newMessage$: Observable<NewMessageResponseBody>;
  public messages$: Observable<any>;
  public deleteMessage$: Observable<DeleteMessageResponseBody>;
  public updateMessage$: Observable<NewMessageResponseBody>;

  constructor(private storageService: StorageService) {
    this.currentUser = this.storageService.userInfo;

    this.newMessage$ = this.messageSubject$.pipe(
      switchAll(),
      filter(message => message.responseType === MessageResponseType.NEW_MESSAGE),
      map(message => NewMessageResponseBody.fromJson(message.body)),
      catchError(e => {
        throw e;
      }),
    );

    this.messages$ = this.messageSubject$.pipe(
      switchAll(),
      filter(message => message.responseType === MessageResponseType.MESSAGE),
      map(message => NewMessageResponseBody.fromJson(message.body)),
      catchError(e => {
        throw e;
      }),
    );

    this.updateMessage$ = this.messageSubject$.pipe(
      switchAll(),
      filter(message => message.responseType === MessageResponseType.UPDATE_MESSAGE),
      map(message => NewMessageResponseBody.fromJson(message.body)),
      catchError(e => {
        throw e;
      }),
    );

    this.deleteMessage$ = this.messageSubject$.pipe(
      switchAll(),
      filter(message => message.responseType === MessageResponseType.DELETE_MESSAGE),
      map(message => DeleteMessageResponseBody.fromJson(message.body)),
      catchError(e => {
        throw e;
      }),
    );
  }

  public connect(cfg: { reconnect: boolean } = {reconnect: false}): void {
    if (!this.socket$ || this.socket$.closed) {
      this.socket$ = this.getNewWebSocket();

      this.registerClient();

      const messages = this.socket$.pipe(
        cfg.reconnect ? this.reconnect : o => o,
        tap({
          error: error => console.error(error)
        }),
        catchError(() => EMPTY),
        filter(data => !!data),
        map(rawData => {
          let listMessageResponse = new Array<MessageResponse<any>>();

          if (_.isObject(rawData)) {
            const rawMessage = MessageResponse.fromJson(rawData);
            listMessageResponse.push(rawMessage);
          } else if (_.isArray(rawData)) {
            listMessageResponse = rawData.map(rawMessage => MessageResponse.fromJson(rawMessage));
          } else {
            const data = '[' + rawData.replace(/\n/g, ',') + '\]';
            const rawMessages = JSON.parse(data);
            listMessageResponse = rawMessages.map(rawMessage => MessageResponse.fromJson(rawMessage));
          }

          return from(listMessageResponse);
        }),
        switchAll()
      );

      this.messageSubject$.next(messages);
    }
  }

  public close(): void {
    this.socket$.complete();
  }

  // region Send Message
  public sendTextMessage(message: string, groupId: number, parentMessageId?: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: SendMessageBody = new SendMessageBody(
      groupId,
      parentMessageId || -1,
      MessageType.TEXT_MESSAGE,
      message
    );

    const request: MessageRequest<SendMessageBody> = new MessageRequest<SendMessageBody>(
      MessageRequestType.SEND_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  public sendImageMessage(imageUrl: string, groupId: number, parentMessageId?: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: SendMessageBody = new SendMessageBody(
      groupId,
      parentMessageId || -1,
      MessageType.IMAGE_MESSAGE,
      imageUrl
    );

    const request: MessageRequest<SendMessageBody> = new MessageRequest<SendMessageBody>(
      MessageRequestType.SEND_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  public sendFileMessage(fileUrl: string, groupId: number, parentMessageId?: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: SendMessageBody = new SendMessageBody(
      groupId,
      parentMessageId || -1,
      MessageType.FILE_MESSAGE,
      fileUrl
    );

    const request: MessageRequest<SendMessageBody> = new MessageRequest<SendMessageBody>(
      MessageRequestType.SEND_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  public updateMessage(groupId: number, messageId: number, content: string, parentMessageId?: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: UpdateMessageBody = new UpdateMessageBody(
      groupId,
      parentMessageId,
      messageId,
      content
    );

    const request: MessageRequest<UpdateMessageBody> = new MessageRequest<UpdateMessageBody>(
      MessageRequestType.UPDATE_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  public deleteMessage(messageId: number, groupId: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: DeleteMessageBody = new DeleteMessageBody(
      groupId,
      -1,
      messageId
    );

    const request: MessageRequest<DeleteMessageBody> = new MessageRequest<DeleteMessageBody>(
      MessageRequestType.DELETE_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  // endregion

  // region Load Message
  public getReplyMessagesHistory(groupId: number, parentMessageId: number, lastMessageId?: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: LoadOldMessageBody = new LoadOldMessageBody(
      groupId,
      parentMessageId,
      lastMessageId
    );

    const request: MessageRequest<LoadOldMessageBody> = new MessageRequest<LoadOldMessageBody>(
      MessageRequestType.LOAD_CHILD_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  public getMessagesHistory(groupId: number, lastMessageId?: number): void {
    const currentUser: User = this.storageService.userInfo;

    const body: LoadOldMessageBody = new LoadOldMessageBody(
      groupId,
      -1,
      lastMessageId
    );

    const request: MessageRequest<LoadOldMessageBody> = new MessageRequest<LoadOldMessageBody>(
      MessageRequestType.LOAD_OLD_MESSAGE,
      currentUser.socketId,
      currentUser.userId,
      body
    );

    this.sendMessage(request);
  }

  // endregion

  private registerClient(): void {
    this.currentUser = this.storageService.userInfo;

    const subscribeMessageBody: SubscribeMessageBody = new SubscribeMessageBody(
      null,
      null,
      this.storageService.token
    );

    _.unset(subscribeMessageBody, 'groupId');
    _.unset(subscribeMessageBody, 'parentMessageId');

    const messageRequest: MessageRequest<SubscribeMessageBody> = new MessageRequest<SubscribeMessageBody>(
      MessageRequestType.SUBSCRIBE,
      this.currentUser.socketId,
      this.currentUser.userId,
      subscribeMessageBody
    );

    this.socket$.next(messageRequest as any);
  }

  private sendMessage(message: MessageRequest<any>): void {
    this.socket$.next(message as any);
  }

  private getNewWebSocket(): WebSocketSubject<MessageResponse<any>> {
    return webSocket({
      url: WS_ENDPOINT,
      serializer: message => JSON.stringify(message),
      deserializer: messageEvent => messageEvent.data,
      closeObserver: {
        next: () => {
          console.log('connection closed');
          this.close();
          this.socket$ = undefined;
          this.connect({reconnect: true});
        }
      }
    });
  }

  private reconnect(observable: Observable<any>): Observable<any> {
    return observable.pipe(
      retryWhen(errors => errors.pipe(
        tap(val => console.log(`Try to reconnect: ${val}`))
      )),
      delay(3000)
    );
  }
}
