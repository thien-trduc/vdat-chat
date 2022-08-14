export class RequestPayload {
  groupId: number;
  userInvite: string;
  userInvited: string;

  constructor(groupId: number, userInvite: string, userInvited: string) {
    this.groupId = groupId;
    this.userInvite = userInvite;
    this.userInvited = userInvited;
  }

  public toJson(): any {
    return {
      idGroup: this.groupId,
      createBy: this.userInvite,
      idInvite: this.userInvited
    };
  }
}
