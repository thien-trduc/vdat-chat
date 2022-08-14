import {Injectable} from '@angular/core';
import {CachingService} from './caching.service';
import * as _ from 'lodash';

@Injectable({
  providedIn: 'root'
})
export class GenerateColorService {

  private colors: Array<{ id: string, color: string }>;

  constructor(private cachingService: CachingService) {
    this.colors = new Array<{id: string; color: string}>();
  }

  public loadCaching(): void {
    this.cachingService.get<Array<{ id: string, color: string }>>(GenerateColorService.name)
      .subscribe(colors => {
        if (!!colors) {
          this.colors = colors;
        }
      });
  }

  public generate(id: string): string {
    const colorFind = _.find(this.colors, iter => iter.id === id);

    if (!!colorFind) {
      return colorFind.color;
    }

    let color = '#';
    const characters = '0123456789ABCDEF';
    const charactersLength = characters.length;

    for (let i = 0; i < 6; i++) {
      color += characters.charAt(Math.floor(Math.random() * charactersLength));
    }

    const existedColor = !!_.find(this.colors, iter => iter.color === color);
    if (existedColor) {
      color = this.generate(id);
    }

    this.colors.push({id, color});

    // save caching
    this.cachingService.save<Array<{ id: string, color: string }>>(GenerateColorService.name, this.colors)
      .subscribe(() => {});

    return color;
  }
}
