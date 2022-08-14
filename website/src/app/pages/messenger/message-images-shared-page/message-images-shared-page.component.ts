import {ChangeDetectionStrategy, Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {FileService} from '../../../core/services/collectors/file.service';
import {ImagesDataSource} from '../../../core/models/datasources/images.datasource';

@Component({
  selector: 'app-message-images-shared-page',
  templateUrl: './message-images-shared-page.component.html',
  styleUrls: ['./message-images-shared-page.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MessageImagesSharedPageComponent implements OnInit, OnChanges {

  @Input() currentGroup: Group;

  public imagesDataSource: ImagesDataSource;

  constructor(private fileService: FileService) {
    this.imagesDataSource = new ImagesDataSource(this.fileService);
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!this.currentGroup) {
      this.imagesDataSource.setCurrentGroup(this.currentGroup.id);
    }
  }
}
