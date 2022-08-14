import {ApplicationRef, Injectable} from '@angular/core';
import {SwUpdate, UpdateAvailableEvent} from '@angular/service-worker';
import {filter, first} from 'rxjs/operators';
import {NzModalService} from 'ng-zorro-antd/modal';
import {concat, interval} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WorkerUpdateService {

  constructor(private appRef: ApplicationRef,
              private swUpdate: SwUpdate,
              private modalService: NzModalService) {
    const appIsStable$ = appRef.isStable.pipe(first(isStable => isStable === true));
    const everySixHours$ = interval(6 * 60 * 60 * 1000);
    const everySixHoursOnceAppIsStable$ = concat(appIsStable$, everySixHours$);

    everySixHoursOnceAppIsStable$.subscribe(() => this.swUpdate.checkForUpdate());

    this.swUpdate.available
      .pipe(filter(event => event.current !== event.available))
      .subscribe(event => {
        console.log('current version is', event.current);
        console.log('available version is', event.available);

        this.updateVersion(event);
      });

    this.swUpdate.activated.subscribe(event => {
      console.log('old version was', event.previous);
      console.log('new version is', event.current);
    });
  }

  private updateVersion(event: UpdateAvailableEvent): void {
    this.modalService.confirm({
      nzTitle: 'Cập nhật phiên bản mới',
      nzContent: `<p>Đã có bản cập nhật mới, vui lòng tải lại trang web hoặc nhấn <b>Cập nhật</b> để cập nhật phiên bản mới !</p><p>Phiên bản hiện tại: ${event.current}</p><p>Phiên bản mới: <b>${event.available}</b></p>`,
      nzOkText: 'Cập nhật ngay',
      nzCancelText: 'Cập nhật sau',
      nzOkType: 'primary',
      nzAutofocus: 'ok',
      nzModalType: 'confirm',
      nzIconType: 'gift',
      nzCentered: true,
      nzOnOk: () => window.location.reload()
    });
  }
}
