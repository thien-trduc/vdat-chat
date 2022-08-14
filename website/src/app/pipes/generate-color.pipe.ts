import {Pipe, PipeTransform} from '@angular/core';
import {GenerateColorService} from '../core/services/commons/generate-color.service';

@Pipe({
  name: 'generateColor'
})
export class GenerateColorPipe implements PipeTransform {

  constructor(private generateColorService: GenerateColorService) {
  }

  transform(key: string): string {
    return this.generateColorService.generate(key);
  }
}
