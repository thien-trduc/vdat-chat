import {CollectionViewer, DataSource} from '@angular/cdk/collections';
import {BehaviorSubject, Observable, Subject} from 'rxjs';

export abstract class BaseDataSource<T> extends DataSource<T> {
  protected pageSize = 10;
  public cachedData: Array<T> = new Array<T>();
  public selecteds: Array<T> = new Array<T>();
  public loading: Subject<boolean> = new BehaviorSubject<boolean>(false);
  protected fetchedPages = new Set<number>();
  protected dataStream = new BehaviorSubject<Array<T>>(this.cachedData);
  protected complete$ = new Subject<void>();
  protected disconnect$ = new Subject<void>();

  protected abstract setup(collectionViewer: CollectionViewer): void;

  public connect(collectionViewer: CollectionViewer): Observable<T[] | ReadonlyArray<T>> {
    this.setup(collectionViewer);
    return this.dataStream;
  }

  public disconnect(collectionViewer: CollectionViewer): void {
    this.disconnect$.next();
    this.disconnect$.complete();
  }

  public completed(): Observable<void> {
    return this.complete$.asObservable();
  }

  public refresh(): void {
    this.cachedData = new Array<T>();
    this.fetchedPages = new Set<number>();
    this.dataStream.next(this.cachedData);
  }

  public toogleLoading(loading: boolean): void {
    this.loading.next(loading);
  }
}
