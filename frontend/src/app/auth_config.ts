import { CustomConfig } from 'ng2-ui-auth';
import { Response } from '@angular/http';

/**
 * Created by Ron on 03/10/2016.
 */
export const GOOGLE_CLIENT_ID = '531545236739-hhtfh2m5rcmeph76sabo3mvdupeu5hfa.apps.googleusercontent.com';

export class MyAuthConfig extends CustomConfig {
    defaultHeaders = { 'Content-Type': 'application/json' };
    providers = {
        google: {
            clientId: GOOGLE_CLIENT_ID,
            url: 'http://localhost:8081/auth/google',
            scope: 'profile email https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/drive.readonly',
            scopeDelimiter: ' ',
            accessType: 'offline'
        }
    };
}