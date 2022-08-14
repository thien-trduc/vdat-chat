import {Component, Input, OnInit} from '@angular/core';
import {Message} from '../../../core/models/message/message.model';
import {NzImage, NzImageService} from "ng-zorro-antd/image";
import {Group} from "../../../core/models/group/group.model";

@Component({
  selector: 'app-message-image',
  templateUrl: './message-image.component.html',
  styleUrls: ['./message-image.component.scss']
})
export class MessageImageComponent implements OnInit {

  @Input() message: Message;
  @Input() currentGroup: Group;
  @Input() isPreviewParentMessage: boolean;

  constructor(private imageService: NzImageService) {
  }

  ngOnInit(): void {
  }

  public onPreviewImage(): void {
    if (!!this.message) {
      const image: NzImage = {
        src: this.message.message,
        alt: `Hình ảnh của ${this.message?.sender?.fullName}`,
      };

      this.imageService.preview([image]);
    }
  }

}
