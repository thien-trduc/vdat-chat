import {BaseDataSource} from './base.datasource';
import * as _ from 'lodash';
import {CollectionViewer} from '@angular/cdk/collections';
import {takeUntil} from 'rxjs/operators';
import {FileService} from '../../services/collectors/file.service';
import {UploadFileDto} from '../file/upload-file.dto';

export class ImagesDataSource extends BaseDataSource<UploadFileDto> {
  private currentGroupId: number;

  constructor(private fileService: FileService) {
    super();
  }

  public setCurrentGroup(currentGroupId: number): void {
    this.currentGroupId = currentGroupId;
    this.refresh();
    this.fetchingData(1);
  }

  protected setup(collectionViewer: CollectionViewer): void {
    this.fetchingData(1);
    collectionViewer.viewChange
      .pipe(
        takeUntil(this.complete$),
        takeUntil(this.disconnect$))
      .subscribe(range => {
      });
  }

  private getPageForIndex(index: number): number {
    return Math.floor(index / this.pageSize);
  }

  public fetchingData(page: number): void {
    if (this.currentGroupId <= 0) {
      return;
    }

    this.cachedData = new Array<UploadFileDto>();
    this.toogleLoading(true);
    this.fetchedPages.add(page);

    this.fileService.getListImageOfGroup(this.currentGroupId)
      .subscribe(image => {
        this.cachedData.push(image);
        this.cachedData = _.uniqBy(this.cachedData, 'nameFile');
        this.dataStream.next(this.cachedData);
        this.toogleLoading(false);
      });
  }
}
