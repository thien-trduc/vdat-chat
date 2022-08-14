import {Pipe, PipeTransform} from '@angular/core';
import {formatDistance, formatDistanceToNow, format} from 'date-fns';
import {vi} from 'date-fns/locale';

@Pipe({
  name: 'datetimeFormat'
})
export class DatetimeFormatPipe implements PipeTransform {
  transform(value: Date = new Date(), formatDate: 'shortTime' | 'relativeTime'): string {
    switch (formatDate) {
      case 'relativeTime':
        return formatDistanceToNow(value, {addSuffix: true, locale: vi});
      case 'shortTime':
        return format(value, 'Pp', {locale: vi});
      default:
        return '';
    }
  }

}
