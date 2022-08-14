import {Inject, Injectable, PLATFORM_ID} from '@angular/core';
import {isPlatformBrowser} from '@angular/common';
import {User} from '../../models/user.model';
import {StorageConst} from '../../constants/storage.const';

@Injectable({
  providedIn: 'root'
})
export class StorageService {

  private readonly isBrowser: boolean;

  constructor(@Inject(PLATFORM_ID) platformId: any) {
    this.isBrowser = isPlatformBrowser(platformId);
  }

  // region User info
  public set userInfo(userInfo: User) {
    if (this.isBrowser) {
      localStorage.setItem(StorageConst.USER_INFO, JSON.stringify(userInfo));
    }
  }

  public get userInfo(): User {
    if (this.isBrowser) {
      const userInfoRaw = localStorage.getItem(StorageConst.USER_INFO);

      if (userInfoRaw) {
        return JSON.parse(userInfoRaw);
      }
    }

    return null;
  }
  // endregion

  // region Token
  public get token(): string {
    return this.isBrowser ? localStorage.getItem(StorageConst.KC_ACCESS_TOKEN) : '';
  }

  public set token(accessToken: string) {
    if (this.isBrowser) {
      localStorage.setItem(StorageConst.KC_ACCESS_TOKEN, accessToken);
    }
  }
  // endregion

  public clearStorage(): void {
    console.log('clear storage');
  }
}
