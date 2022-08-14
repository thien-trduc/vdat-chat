import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DelayInputDirective} from "./delay-input.directive";


@NgModule({
  declarations: [
    DelayInputDirective
  ],
  exports: [
    DelayInputDirective
  ],
  imports: [
    CommonModule
  ]
})
export class CustomDirectiveModule {
}
