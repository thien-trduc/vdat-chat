import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserInfoPipe} from './user-info.pipe';
import {MessagePipe} from './message.pipe';
import {GroupInfoPipe} from './group/group-info.pipe';
import {GroupRolePipe} from './group/group-role.pipe';
import {DatetimeFormatPipe} from './datetime-format.pipe';
import {FileInfoPipe} from './file-info.pipe';
import {GenerateColorPipe} from './generate-color.pipe';

@NgModule({
  declarations: [
    UserInfoPipe,
    MessagePipe,
    GroupInfoPipe,
    GroupRolePipe,
    DatetimeFormatPipe,
    FileInfoPipe,
    GenerateColorPipe,
  ],
  exports: [
    UserInfoPipe,
    MessagePipe,
    GroupRolePipe,
    DatetimeFormatPipe,
    GroupInfoPipe,
    FileInfoPipe,
  ],
  imports: [
    CommonModule
  ]
})
export class CustomPipesModule {
}
