// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
  redirectUrl: 'http://localhost:4200',
  googleClientId: '531545236739-hhtfh2m5rcmeph76sabo3mvdupeu5hfa.apps.googleusercontent.com',
  slackClientId: '16724837808.155634746450',
  githubClientId: '469a838ef4c6048510b6'
};
