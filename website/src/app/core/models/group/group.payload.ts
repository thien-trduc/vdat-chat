import {GroupType} from '../../constants/group-type.const';

export interface GroupPayload {
  nameGroup: string;
  type: GroupType;
  private: boolean;
  users: Array<string>;
  description: string;
  thumbnail?: string;
}
