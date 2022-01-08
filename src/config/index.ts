import { Command } from 'commander';
import { getConfig } from './funcs';

const infoCmd = new Command('info').action(() => {
    const config = getConfig();
    console.log(`
Soar Config
====================
Application Details

url: ${config.application.url}
key: ${config.application.key}
uses: 0

Client Details
url: ${config.client.url}
key: ${config.client.key}

General
use strict mode:            ${config.logs.strictMode}
output debug messages:      ${config.logs.showDebug}
output http responses:      ${config.logs.showHTTPLog}
output websocket responses: ${config.logs.showWebSocketLog}

Error Directory: ${config.logs.errorOutDir}
    `);
});

const main = new Command('config')
    .addHelpText('before', 'Manages the internal Soar configurations.')
    .addCommand(infoCmd);

export default main;
