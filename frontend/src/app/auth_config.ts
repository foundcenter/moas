import {CustomConfig} from 'ng2-ui-auth';
/**
 * Created by Ron on 03/10/2016.
 */
export const GOOGLE_CLIENT_ID = '531545236739-hhtfh2m5rcmeph76sabo3mvdupeu5hfa.apps.googleusercontent.com';
export const SLACK_CLIENT_ID = '16724837808.155634746450';

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
            url: "https://httpbin.org/post",
            authorizationEndpoint: 'https://slack.com/oauth/authorize',
            requiredUrlParams: ['scope'],
            redirectUri: 'http://localhost:4200',
            scope: 'files:read mpim:read search:read',
            scopePrefix: '',
            scopeDelimiter: ',',
            oauthType: '2.0',
            popupOptions: { width: 700, height: 800 },
        }
    };
}