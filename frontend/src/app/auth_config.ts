import {CustomConfig} from 'ng2-ui-auth';
/**
 * Created by Ron on 03/10/2016.
 */
export const GOOGLE_CLIENT_ID = '531545236739-hhtfh2m5rcmeph76sabo3mvdupeu5hfa.apps.googleusercontent.com';
export const SLACK_CLIENT_ID = '16724837808.155634746450';
export const GITHUB_CLIENT_ID = '469a838ef4c6048510b6';

export const BACKEND_MOCK = 'https://httpbin.org/post';

export class MyAuthConfig extends CustomConfig {
    defaultHeaders = {'Content-Type': 'application/json'};
    providers = {
        google: {
            clientId: GOOGLE_CLIENT_ID,
            url: "http://localhost:8081/auth/google",
            scope: "profile email https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/drive.readonly",
            scopeDelimiter: " "
        },
        slack: {
            clientId: SLACK_CLIENT_ID,
            url: BACKEND_MOCK,
            authorizationEndpoint: 'https://slack.com/oauth/authorize',
            requiredUrlParams: ['scope'],
            redirectUri: 'http://localhost:4200',
            scope: 'files:read mpim:read search:read',
            scopePrefix: '',
            scopeDelimiter: ',',
            oauthType: '2.0',
            popupOptions: { width: 700, height: 800 },
        },
        github: {
            clientId: GITHUB_CLIENT_ID,
            name: 'github',
            url: BACKEND_MOCK,
            authorizationEndpoint: 'https://github.com/login/oauth/authorize',
            redirectUri: 'http://localhost:4200',
            requiredUrlParams: ['scope'],
            scope: ['user:email'],
            scopeDelimiter: ' ',
            oauthType: '2.0',
            popupOptions: {width: 1020, height: 720}
        },
    };
}