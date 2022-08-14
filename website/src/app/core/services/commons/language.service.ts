import { Injectable } from '@angular/core';
import {TranslateService} from '@ngx-translate/core';

@Injectable({
  providedIn: 'root'
})
export class LanguageService {
  private currentLanguage: string;

  constructor(private translateService: TranslateService) {
    this.currentLanguage = 'vi';
  }

  public getCurrentLanguage(): string {
    return this.currentLanguage;
  }

  public async translation(text: string): Promise<string> {
    return await this.translateService.get(text).toPromise();
  }

  public setDefaultLanguage(): void {
    this.translateService.setDefaultLang(this.currentLanguage);
    this.switchLanguage(this.currentLanguage);
  }

  public switchLanguage(language: string): void {
    this.currentLanguage = language;
    this.translateService.use(this.currentLanguage);
  }
}
