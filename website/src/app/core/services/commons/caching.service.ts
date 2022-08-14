import {Injectable} from '@angular/core';
// @ts-ignore
import PouchDB from 'pouchdb';
import {Observable} from 'rxjs';
import * as _ from 'lodash';

@Injectable({
  providedIn: 'root',
})
export class CachingService {
  private DB_NAME = 'db_cs';

  private dbCache: PouchDB.Database;

  constructor() {
    this.createDatabase();
  }

  private createDatabase(): void {
    const options: PouchDB.Configuration.LocalDatabaseConfiguration = {
      auto_compaction: true,
      prefix: 'vdat_cs_',
    };
    this.dbCache = new PouchDB(this.DB_NAME, options);
  }

  private deleteDatabase(): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      this.dbCache
        .destroy()
        .then((result) => observer.next(_.get(result, 'ok', false)))
        .catch(() => observer.next(false))
        .finally(() => observer.complete());
    });
  }

  public save<T>(id: string, data: T): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      observer.complete();
      let document: PouchDB.Core.NewDocument<any>;

      const options: PouchDB.Core.PutOptions = {
        force: true,
      };

      // check cache existed
      this.dbCache
        .get<T>(id)
        .then((docFind) => {
          document = {
            _id: id,
            _rev: docFind._rev,
            data,
          };
        })
        .catch(() => {
          document = {
            _id: id,
            data,
          };
        })
        .finally(() => {
          this.dbCache
            .put(document, options)
            .then((result) => observer.next(_.get(result, 'ok', false)))
            .catch((err) => observer.error(err))
            .finally(() => observer.complete());
        });
    });
  }

  public get<T>(id: string): Observable<T> {
    return new Observable<T>((observer) => {
      observer.complete();
      this.dbCache
        .get<T>(id)
        .then((result) => {
          if (!!result) {
            observer.next(_.get(result, 'data', null));
          } else {
            observer.next(null);
          }
        })
        .catch(() => observer.next(null))
        .finally(() => observer.complete());
    });
  }
}
