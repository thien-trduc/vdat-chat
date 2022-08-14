import { Injectable } from '@angular/core';
import * as CryptoJS from 'crypto-js';
import {environment} from '../../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class EncryptService {

  private secretKey = environment.secretKey;

  constructor() { }

  public encrypt(data: string): string {
    const wordArray = CryptoJS.enc.Utf8.parse(data);
    const base64 = CryptoJS.enc.Base64.stringify(wordArray);
    return base64.toString();
  }

  public decrypt(cipherText: string): string {
    const parsedWordArray = CryptoJS.enc.Base64.parse(cipherText);
    const parsedStr = parsedWordArray.toString(CryptoJS.enc.Utf8);
    return parsedStr.toString();
  }
}
