const IOH_SERVER_URL = 'vdat-mcsvc-collector-staging.vdatlab.com';

export const environment = {
  production: true,
  api: {
    users: '/api/v1/user',
    groups: '/api/v1/groups',
    files: '/api/v1/files',
    request: '/api/v1/requests',
    messages: {
      path: '/message',
      protocol: 'wss'
    }
  },
  ioh: {
    apiUrl: `https://${IOH_SERVER_URL}/dcs/v1`,
    endpoint: {
      user: 'users'
    }
  },
  keycloak: {
    url: 'https://accounts.vdatlab.com/auth',
    realm: 'vdatlab.com',
    clientId: 'chat.apps.vdatlab.com'
  },
  secretKey: 'test123'
};
