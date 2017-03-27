import { CustomConfig } from 'ng2-ui-auth';
import { Response } from '@angular/http';
import { environment } from '../environments/environment';

/**
 * Created by Ron on 03/10/2016.
 */
export const GOOGLE_CLIENT_ID = environment.googleClientId;
export const SLACK_CLIENT_ID = environment.slackClientId;
export const GITHUB_CLIENT_ID = environment.githubClientId;

export const BACKEND_MOCK = 'https://httpbin.org/post';
export const REDIRECT_URI = environment.redirectUrl;
export const API_URL = environment.apiUrl;

export class MyAuthConfig extends CustomConfig {
    defaultHeaders = {'Content-Type': 'application/json'};
    providers = {
        google: {
            clientId: GOOGLE_CLIENT_ID,
            url: `${API_URL}/auth/google`,
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
            url: `${API_URL}/connect/drive`,
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
            scope: 'search:read identify users:read',
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
            redirectUri: REDIRECT_URI,
            requiredUrlParams: ['scope'],
            scope: 'user:email',
            scopeDelimiter: ' ',
            oauthType: '2.0',
            popupOptions: {width: 1020, height: 720}
        }
    };
}