import {Component, HostListener} from '@angular/core';
import {LanguageService} from './core/services/commons/language.service';
import {GenerateColorService} from './core/services/commons/generate-color.service';
import {UserService} from './core/services/collectors/user.service';
import {StorageService} from './core/services/commons/storage.service';
import {MessageService} from './core/services/ws/message.service';
import {retry} from "rxjs/operators";
import {NzMessageService} from "ng-zorro-antd/message";
import {WorkerUpdateService} from "./core/services/workers/worker-update.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  public loading: boolean;

  constructor(private languageService: LanguageService,
              private userService: UserService,
              private storageService: StorageService,
              private messageService: MessageService,
              private nzMessageService: NzMessageService,
              private generateColorService: GenerateColorService) {
    this.loading = true;
    this.languageService.setDefaultLanguage();
    this.generateColorService.loadCaching();

    this.userService.getUserInfo()
      .subscribe(user => this.storageService.userInfo = user);
  }

  @HostListener('window:beforeunload', ['$event'])
  public beforeUnloadHandler(event): void {
    this.userService.logout().subscribe(() => {
      event.returnValue = true;
    });
  }
}
