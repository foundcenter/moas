import { CustomConfig } from 'ng2-ui-auth';
import { Response } from '@angular/http';

/**
 * Created by Ron on 03/10/2016.
 */
export const GOOGLE_CLIENT_ID = '531545236739-hhtfh2m5rcmeph76sabo3mvdupeu5hfa.apps.googleusercontent.com';
export const SLACK_CLIENT_ID = '16724837808.155634746450';
export const GITHUB_CLIENT_ID = '469a838ef4c6048510b6';

export const BACKEND_MOCK = 'https://httpbin.org/post';
export const REDIRECT_URI = 'http://localhost:4200';
export const API_URL ='http://localhost:8081';

export class MyAuthConfig extends CustomConfig {
    defaultHeaders = {'Content-Type': 'application/json'};
    providers = {
        google: {
            clientId: GOOGLE_CLIENT_ID,
            url: 'http://localhost:8081/auth/google',
            scope: 'profile email https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/drive.readonly',
            scopeDelimiter: ' ',
            accessType: 'offline'
        },
        gmail: {
            name: 'gmail',
            clientId: GOOGLE_CLIENT_ID,
            url: `${API_URL}/connect/gmail`,
            authorizationEndpoint: 'https://accounts.google.com/o/oauth2/auth',
            redirectUri: REDIRECT_URI,
            requiredUrlParams: ['scope'],
            optionalUrlParams: ['display', 'state', 'prompt', 'login_hint', 'access_type',
                'include_granted_scopes', 'openid.realm', 'hd'],
            scope: 'profile email https://www.googleapis.com/auth/gmail.readonly',
            scopePrefix: 'openid',
            scopeDelimiter: ' ',
            accessType: 'offline',
            display: 'popup',
            oauthType: '2.0',
            popupOptions: {width: 452, height: 633},
            state: () => encodeURIComponent(Math.random().toString(36).substr(2)),
        },
        'google-drive': {
            name: 'google-drive',
            clientId: GOOGLE_CLIENT_ID,
            url: BACKEND_MOCK,
            authorizationEndpoint: 'https://accounts.google.com/o/oauth2/auth',
            redirectUri: REDIRECT_URI,
            requiredUrlParams: ['scope'],
            optionalUrlParams: ['display', 'state', 'prompt', 'login_hint', 'access_type',
                'include_granted_scopes', 'openid.realm', 'hd'],
            scope: 'profile email https://www.googleapis.com/auth/drive.readonly',
            scopePrefix: 'openid',
            scopeDelimiter: ' ',
            accessType: 'offline',
            display: 'popup',
            oauthType: '2.0',
            popupOptions: {width: 452, height: 633},
            state: () => encodeURIComponent(Math.random().toString(36).substr(2)),
        },
        slack: {
            clientId: SLACK_CLIENT_ID,
            url:  `${API_URL}/connect/slack`,
            authorizationEndpoint: 'https://slack.com/oauth/authorize',
            requiredUrlParams: ['scope'],
            redirectUri: REDIRECT_URI,
            scope: 'search:read identity.basic identity.email',
            scopePrefix: '',
            scopeDelimiter: ',',
            oauthType: '2.0',
            popupOptions: { width: 700, height: 800 },
        },
        github: {
            clientId: GITHUB_CLIENT_ID,
            name: 'github',
            url: `${API_URL}/connect/github`,
            authorizationEndpoint: 'https://github.com/login/oauth/authorize',
            redirectUri: 'http://localhost:4200',
            requiredUrlParams: ['scope'],
            scope: 'user:email',
            scopeDelimiter: ' ',
            oauthType: '2.0',
            popupOptions: {width: 1020, height: 720}
        }
    };
}