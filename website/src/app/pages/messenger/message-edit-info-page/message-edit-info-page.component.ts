import {Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {GroupPayload} from '../../../core/models/group/group.payload';
import * as _ from 'lodash';

@Component({
  selector: 'app-message-edit-info-page',
  templateUrl: './message-edit-info-page.component.html',
  styleUrls: ['./message-edit-info-page.component.scss']
})
export class MessageEditInfoPageComponent implements OnInit, OnChanges {

  @Input() currentGroup: Group;

  public formGroup: FormGroup;
  public loading: boolean;
  public uploadApiEndpoint: string;
  public nameFileUploaded: string;

  constructor() {
    this.formGroup = this.createFormGroup();
    this.uploadApiEndpoint = `/api/v1/files/upload`;
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && changes.currentGroup && !!this.currentGroup) {
      this.formGroup = this.createFormGroup(this.currentGroup);
    }
  }

  ngOnInit(): void {
  }

  public getFormValue(): GroupPayload {
    if (this.formGroup.valid) {
      const rawValue = this.formGroup.getRawValue();

      if (!!rawValue) {
        const isPrivate = _.get(rawValue, 'private', true);
        _.set(rawValue, 'private', !isPrivate);

        if (!!this.nameFileUploaded) {
          _.set(rawValue, 'thumbnail', this.nameFileUploaded);
        }

        return rawValue;
      }
    }

    return null;
  }

  private createFormGroup(group?: Group): FormGroup {
    const formGroup: FormGroup = new FormGroup({
      nameGroup: new FormControl('', Validators.required),
      description: new FormControl(''),
      private: new FormControl(false),
      thumbnail: new FormControl('')
    });

    if (!!group) {
      _.set(group, 'private', !group.isPrivate);
      formGroup.patchValue(group);
    }

    return formGroup;
  }
}
