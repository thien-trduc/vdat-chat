import {Injectable} from '@angular/core';
import {KeycloakInstance, KeycloakLoginOptions, KeycloakInitOptions, KeycloakLogoutOptions} from 'keycloak-js';
import {StorageConst} from '../../constants/storage.const';
import * as Keycloak from 'keycloak-js';
import {Observable, Subject} from 'rxjs';
import {environment} from '../../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class KeycloakService {
  private readonly ACCESS_TOKEN = StorageConst.KC_ACCESS_TOKEN;
  private readonly REFRESH_TOKEN = StorageConst.KC_REFRESH_TOKEN;
  private readonly ID_TOKEN = StorageConst.KC_ID_TOKEN;
  private readonly USER_INFO = StorageConst.KC_USER_INFO;
  private readonly TIME_REFRESH = 120;

  private keycloak: KeycloakInstance;
  public tokenListener: Subject<string> = new Subject<string>();

  public getKeycloakInstance(clientId?: string, integrated: boolean = false): Observable<KeycloakInstance> {
    return new Observable<Keycloak.KeycloakInstance>(observer => {
      if (!!this.keycloak) {
        observer.next(this.keycloak);
        observer.complete();
        return;
      }

      this.keycloak = Keycloak(this.getKeycloakConfig(clientId));

      const initOptions: KeycloakInitOptions = {
        onLoad: integrated ? 'login-required' : 'check-sso',
        checkLoginIframe: !integrated,
        checkLoginIframeInterval: this.TIME_REFRESH,
        idToken: this.idToken,
        token: this.accessToken,
        refreshToken: this.refreshToken,
        redirectUri: `${window.location.origin}/auth`
      };

      this.keycloak.init(initOptions)
        .then(() => {
          setTimeout(() => {
            this.keycloak.updateToken(600)
              .then((refreshed) => {
                if (refreshed) {
                  console.log('Token refreshed ' + refreshed);
                } else {
                  console.warn('Token not refreshed, valid for '
                    + Math.round(this.keycloak.tokenParsed.exp + this.keycloak.timeSkew - new Date().getTime() / 1000) + ' seconds');
                }
              }).catch(() => {
              console.error('Failed to refresh token');
            });
          }, 60000);

          observer.next(this.keycloak);
          observer.complete();
        })
        .catch(err => {
          console.log(err);
          this.clearAuth();
          observer.next(null);
          observer.complete();
        });

      this.keycloak.onReady = this.onReady;
      this.keycloak.onAuthSuccess = this.onAuthSuccess;
      this.keycloak.onAuthError = this.onAuthError;
      this.keycloak.onAuthRefreshSuccess = this.onAuthSuccess;
      this.keycloak.onAuthRefreshError = this.onAuthError;
    });
  }

  public onReady = (authenticated: boolean) => {
    if (authenticated) {
      this.onAuthSuccess();
    }
  }

  public onAuthSuccess = () => {
    this.accessToken = this.keycloak.token;
    this.refreshToken = this.keycloak.refreshToken;
    this.idToken = this.keycloak.idToken;

    this.keycloak.loadUserInfo()
      .then(userInfo => {
        this.userInfo = userInfo;
      });

    this.tokenListener.next(this.accessToken);
  }

  public onAuthError = () => {
    this.clearAuth();

    this.tokenListener.next(null);
  }

  public forceCreate(clientId?: string, integrated: boolean = false): Observable<KeycloakInstance> {
    this.keycloak = null;
    return this.getKeycloakInstance(clientId, integrated);
  }

  public login(options?: KeycloakLoginOptions): void {
    this.getKeycloakInstance()
      .subscribe(keycloak => {
        keycloak.login(options);
      });
  }

  public logout(options?: KeycloakLogoutOptions): void {
    this.getKeycloakInstance()
      .subscribe(keycloak => {
        keycloak.logout(options);
        this.clearAuth();
      });
  }

  private getKeycloakConfig(clientId?: string): any {
    return {
      url: environment.keycloak.url,
      realm: environment.keycloak.realm,
      clientId: clientId || environment.keycloak.clientId
    };
  }

  // region Storage
  public set userInfo(userInfo: any) {
    localStorage.setItem(this.USER_INFO, JSON.stringify(userInfo));
  }

  public get userInfo(): any {
    const userInfoRaw = localStorage.getItem(this.USER_INFO);

    if (userInfoRaw) {
      return JSON.parse(userInfoRaw);
    }
    return null;
  }

  public set idToken(idToken: string) {
    localStorage.setItem(this.ID_TOKEN, idToken);
  }

  public get idToken(): string {
    return localStorage.getItem(this.ID_TOKEN);
  }

  public set refreshToken(refreshToken: string) {
    localStorage.setItem(this.REFRESH_TOKEN, refreshToken);
  }

  public get refreshToken(): string {
    return localStorage.getItem(this.REFRESH_TOKEN);
  }

  public set accessToken(accessToken: string) {
    localStorage.setItem(this.ACCESS_TOKEN, accessToken);
  }

  public get accessToken(): string {
    return localStorage.getItem(this.ACCESS_TOKEN);
  }

  public clearAuth(): void {
    localStorage.removeItem(this.ACCESS_TOKEN);
    localStorage.removeItem(this.REFRESH_TOKEN);
    localStorage.removeItem(this.USER_INFO);
  }

  // endregion
}
