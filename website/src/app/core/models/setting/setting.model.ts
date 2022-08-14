import {SoundSetting} from './sound-setting.model';

export class Setting {
  sound: SoundSetting;
  isDarkMode: boolean;
  isNotification: boolean;

  constructor(sound: SoundSetting = null, isDarkMode: boolean = false, isNotification: boolean = false) {
    this.sound = sound;
    this.isDarkMode = isDarkMode;
    this.isNotification = isNotification;
  }
}
