import * as _ from 'lodash';

export class UploadFileDto {
  location: string;
  nameFile: string;
  fileUrl: string;
  type: string;

  constructor(location: string, nameFile: string, fileUrl: string, type: string) {
    this.location = location;
    this.nameFile = nameFile;
    this.fileUrl = fileUrl;
    this.type = type;
  }

  public static fromJson(obj: any): UploadFileDto {
    return new UploadFileDto(
      _.get(obj, 'Location', ''),
      _.get(obj, 'NameFile', ''),
      _.get(obj, 'ShareUrl', ''),
      _.get(obj, 'Type', '')
    );
  }
}
