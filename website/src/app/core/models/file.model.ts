export class FileModel {
  fileName: string;
  fileType: string;
  fileUrl: string;
  icon: string;

  constructor(fileName: string, fileUrl: string) {
    this.fileName = fileName;
    this.fileUrl = fileUrl;
    this.fileType = this.getFileType();
    this.icon = this.getIconByType();
  }

  private getIconByType(): string {
    switch (this.fileType) {
      case 'gif':
        return 'file-gif';
      case 'md':
        return 'file-markdown';
      case 'pdf':
        return 'file-pdf';
      case 'txt':
        return 'file-text';
      case 'ppt':
      case 'pptx':
      case 'odp':
        return 'file-ppt';
      case 'css':
      case 'csv':
      case 'ods':
      case 'xls':
      case 'xlsx':
        return 'file-excel';
      case 'doc':
      case 'docx':
      case 'odt':
        return 'file-word';
      case 'zip':
      case 'rar':
      case '7z':
      case 'gz':
      case 'bz':
      case 'bz2':
      case 'tar':
        return 'file-zip';
      default:
        return 'file-unknown';
    }
  }

  private getFileType(): string {
    return (/[.]/.exec(this.fileName)) ? /[^.]+$/.exec(this.fileName)[0] : undefined;
  }
}
