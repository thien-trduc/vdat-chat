export class SelectItem<T> {
  data: T;
  checked: boolean;

  constructor(data: T, checked = false) {
    this.data = data;
    this.checked = checked;
  }

  public select(): void {
    this.checked = true;
  }

  public unSelect(): void {
    this.checked = false;
  }

  public toggleSelect(): void {
    this.checked = !this.checked;
  }
}
