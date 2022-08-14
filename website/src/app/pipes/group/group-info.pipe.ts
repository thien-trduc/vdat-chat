import { Pipe, PipeTransform } from '@angular/core';

const argsType = 'is-owner';

@Pipe({
  name: 'groupInfo'
})
export class GroupInfoPipe implements PipeTransform {
  transform(value: unknown, args: 'is-owner' | 'is-member'): boolean {
    return null;
  }
}
