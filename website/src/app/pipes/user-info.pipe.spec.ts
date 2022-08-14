import { UserInfoPipe } from './user-info.pipe';

describe('UserInfoPipe', () => {
  it('create an instance', () => {
    const pipe = new UserInfoPipe();
    expect(pipe).toBeTruthy();
  });
});
