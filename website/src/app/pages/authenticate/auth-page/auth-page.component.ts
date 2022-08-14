import {Component, Inject, OnInit} from '@angular/core';
import {DOCUMENT} from "@angular/common";
import {KeycloakService} from "../../../core/services/auth/keycloak.service";

@Component({
  selector: 'app-auth-page',
  templateUrl: './auth-page.component.html',
  styleUrls: ['./auth-page.component.scss']
})
export class AuthPageComponent implements OnInit {

  constructor(@Inject(DOCUMENT) private document: Document,
              private keycloakService: KeycloakService) {
    this.keycloakService.getKeycloakInstance()
      .subscribe(keycloak => {
        if (keycloak.authenticated) {
          document.location.href = '/';
        }
      });
  }

  ngOnInit(): void {
  }

  public onLogin(): void {
    this.keycloakService.login();
  }

}
