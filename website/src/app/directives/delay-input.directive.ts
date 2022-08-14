import {Directive, Input, OnDestroy, OnInit, Output, EventEmitter, ElementRef} from '@angular/core';
import {fromEvent, Subject, timer} from 'rxjs';
import {debounce, distinctUntilChanged, filter, map, takeUntil} from 'rxjs/operators';

@Directive({
  selector: '[appDelayInput]'
})
export class DelayInputDirective implements OnInit, OnDestroy {
  private destroy$ = new Subject<void>();

  @Input() delayTime = 500;
  @Output() delayedInput = new EventEmitter<Event>();

  constructor(private elementRef: ElementRef<HTMLInputElement>) {
  }

  ngOnInit(): void {
    fromEvent(this.elementRef.nativeElement, 'input')
      .pipe(
        debounce(() => timer(this.delayTime)),
        distinctUntilChanged(
          null,
          (event: Event) => (event.target as HTMLInputElement).value
        ),
        takeUntil(this.destroy$),
      )
      .subscribe(event => this.delayedInput.emit(event));
  }

  ngOnDestroy(): void {
    this.destroy$.next();
  }
}
