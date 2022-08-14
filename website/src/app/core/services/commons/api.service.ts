import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpHeaders, HttpParams, HttpResponse} from '@angular/common/http';
import {StorageService} from './storage.service';
import * as _ from 'lodash';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private httpClient: HttpClient,
              private storageService: StorageService) {
  }

  public setHeaders(headers?: any): HttpHeaders {
    const token = 'Bearer ' + this.storageService.token;
    let httpHeaders;

    if (token) {
      try {
        httpHeaders = new HttpHeaders(_.assign({
          'Content-Type': 'application/json; charset=utf-8',
          Authorization: token
        }, headers));
      } catch (error) {
        this.storageService.clearStorage();
      }
    }

    return httpHeaders;
  }

  public setUploadFileHeaders(headers?: any): HttpHeaders {
    const token = 'Bearer ' + this.storageService.token;
    let httpHeaders;

    if (token) {
      try {
        httpHeaders = new HttpHeaders(_.assign({
          Authorization: token
        }, headers));
      } catch (error) {
        this.storageService.clearStorage();
      }
    }

    return httpHeaders;
  }

  public setUrlEncodedHeaders(headers?: any): HttpHeaders {
    const token = 'Bearer ' + this.storageService.token;
    let httpHeaders;

    if (token) {
      try {
        httpHeaders = new HttpHeaders(_.assign({
          'Content-Type': 'application/x-www-form-urlencoded',
          Authorization: token
        }, headers));
      } catch (error) {
        this.storageService.clearStorage();
      }
    }

    return httpHeaders;
  }

  private errorHandler(error: HttpErrorResponse): Observable<never> {
    if (error.error instanceof ErrorEvent) {
      console.error('An error occurred: ', error.error.message);
    } else {
      return throwError(error);
    }
    return throwError('Something went wrong!');
  }

  public post(path: string, body: any, customHeader?: any): Observable<HttpResponse<any>> {
    return this.httpClient.post<any>(
      path, body,
      {
        headers: this.setHeaders(customHeader),
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public postUrlEncoded(path: string, body: any, customHeader?: any): Observable<HttpResponse<any>> {
    return this.httpClient.post<any>(
      path, body,
      {
        headers: this.setUrlEncodedHeaders(customHeader),
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public postFile(path: string, body: any, customHeader?: any): Observable<HttpResponse<any>> {
    return this.httpClient.post<any>(
      path, body,
      {
        headers: this.setUploadFileHeaders(customHeader),
        reportProgress: true,
        withCredentials: false,
        observe: 'response',
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public get(path: string, options?, params?: HttpParams): Observable<any> {
    return this.httpClient.get(
      path,
      {
        headers: this.setHeaders(options),
        params,
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public put(path: string, body?: any): Observable<any> {
    return this.httpClient.put(
      path, body,
      {
        headers: this.setHeaders(),
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public patch(path: string, body?: any): Observable<any> {
    return this.httpClient.patch(
      path, body,
      {
        headers: this.setHeaders(),
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public patchFile(path: string, body?: any): Observable<any> {
    return this.httpClient.patch(
      path, body,
      {
        headers: this.setUploadFileHeaders(),
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }

  public delete(path: string): Observable<any> {
    return this.httpClient.delete(
      path,
      {
        headers: this.setHeaders(),
        withCredentials: false,
        observe: 'response'
      })
      .pipe(
        catchError(this.errorHandler)
      );
  }
}
