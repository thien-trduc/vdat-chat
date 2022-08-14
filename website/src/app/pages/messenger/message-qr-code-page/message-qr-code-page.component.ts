import {Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import QRCode from 'qrcode';
import {User} from '../../../core/models/user.model';
import {EncryptService} from '../../../core/services/commons/encrypt.service';
import {Invite} from '../../../core/models/invite.model';

@Component({
  selector: 'app-message-qr-code-page',
  templateUrl: './message-qr-code-page.component.html',
  styleUrls: ['./message-qr-code-page.component.scss']
})
export class MessageQrCodePageComponent implements OnInit, OnChanges {

  @Input() currentGroup: Group;
  @Input() currentUser: User;

  public generateQRCode: string;
  public linkInvite: string;

  constructor(private encryptService: EncryptService) {
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!this.currentGroup) {
      const groupId = this.currentGroup?.id;
      const groupName = this.currentGroup?.nameGroup;
      const userInviteId = this.currentUser?.userId;
      const userInviteFullName = this.currentUser?.fullName;
      const timeInvite = new Date();

      const invite: Invite = new Invite(groupId, groupName, userInviteId, userInviteFullName, timeInvite);

      const data = invite.toData(this.encryptService);
      this.linkInvite = `${window.location.protocol}//${window.location.host}/invite?g=${data}`;

      QRCode.toDataURL(this.linkInvite, this.getOptionsGenerateQrCode())
        .then(url => this.generateQRCode = url)
        .catch(err => this.generateQRCode = null);
    }
  }

  public downloadQrCode(): void {
    console.log('download qr code');
  }

  private getOptionsGenerateQrCode(): any {
    return {
      margin: 4,
      scale: 4,
      width: 256,
      errorCorrectionLevel: 'H',
      maskPattern: 7
    };
  }
}

