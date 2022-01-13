import { Command } from 'commander';
import { getConfig } from './funcs';

import setupCmd from './setup';

const infoCmd = new Command('info').action(() => {
    const config = getConfig();
    console.log(`
Soar Config
====================
Application Details

url: ${config.application.url || 'Not Set'}
key: ${config.application.key || 'Not Set'}
uses: 0

Client Details
url: ${config.client.url || 'Not Set'}
key: ${config.client.key || 'Not Set'}

General
use strict mode:            ${config.logs.strictMode}
output debug messages:      ${config.logs.showDebug}
output http responses:      ${config.logs.showHttpLog}
output websocket responses: ${config.logs.showWsLog}

Error Directory: ${config.logs.errorOutDir || 'Not Set'}
    `);
});

const main = new Command('config')
    .addHelpText('before', 'Manages the internal Soar configurations.')
    .addCommand(infoCmd)
    .addCommand(setupCmd);

export default main;
