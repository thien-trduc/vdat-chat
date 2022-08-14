import {Injectable} from '@angular/core';
import {CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router, CanActivateChild} from '@angular/router';
import {Observable} from 'rxjs';
import {KeycloakService} from '../core/services/auth/keycloak.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate, CanActivateChild {
  constructor(private keycloakService: KeycloakService,
              private router: Router) {
  }

  canActivate(next: ActivatedRouteSnapshot,
              state: RouterStateSnapshot)
    : Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return this.authentication();
  }

  canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot)
    : Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return this.authentication();
  }

  private authentication(): Observable<boolean> {
    return new Observable<boolean>(observer => {
      this.keycloakService.getKeycloakInstance()
        .subscribe(keycloak => {
          if (!!keycloak && keycloak.authenticated) {
            observer.next(true);
            observer.complete();
          } else {
            observer.next(false);
            this.router.navigateByUrl('/auth').then(() => observer.complete());
          }
        });
    });
  }
}
