import {Pipe, PipeTransform} from '@angular/core';
import {from, Observable} from 'rxjs';
import {filter, map} from 'rxjs/operators';

@Pipe({
  name: 'fileInfo'
})
export class FileInfoPipe implements PipeTransform {

  transform(fileName: string, key: 'fileName' | 'iconFile'): any {
    switch (key) {
      case 'fileName':
        return this.getFileName(fileName);
      case 'iconFile':
        return this.getIconFile(fileName);
    }
  }

  private getFileName(rawFileName: string): string {
    const startIndex = Math.max(rawFileName.indexOf('_', 0), 0);
    return rawFileName.substring(startIndex > 0 ? startIndex + 1 : startIndex);
  }

  private getIconFile(rawFileName: string): Observable<string> {
    const fileName = this.getFileName(rawFileName);
    const extFile = fileName.split('.').pop().trim();

    return from(this.getListFileIcon())
      .pipe(filter(obj => !!obj))
      .pipe(filter(obj => !!obj.extensions.split('|').find(ext => ext === extFile)))
      .pipe(map(obj => obj.icon));
  }

  private getListFileIcon(): Array<{extensions: string, icon: string}> {
    return [
      {
        extensions: 'txt|tex',
        icon: 'file-text'
      },
      {
        extensions: 'gif',
        icon: 'file-gif'
      },
      {
        extensions: 'md',
        icon: 'file-markdown'
      },
      {
        extensions: 'pdf',
        icon: 'file-pdf'
      },
      {
        extensions: 'ai|bmp|ico|jpeg|jpg|png|ps|psd|svg|tif|tiff',
        icon: 'file-image'
      },
      {
        extensions: 'doc|docx|odt|rtf|wpd',
        icon: 'file-word'
      },
      {
        extensions: 'ods|xls|xlsm|xlsx',
        icon: 'file-excel'
      },
      {
        extensions: 'key|odp|ppt|pptx',
        icon: 'file-ppt'
      },
      {
        extensions: '7z|arj|rar|gz|z|zip',
        icon: 'file-zip'
      },
      {
        extensions: '*',
        icon: 'file'
      }
    ];
  }
}
