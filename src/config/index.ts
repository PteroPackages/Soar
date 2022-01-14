import { Command } from 'commander';
import { getConfig } from './funcs';

import setupCmd from './setup';

const infoCmd = new Command('info')
    .option('-h, --hide', 'Hides the API keys from the command output.', false)
    .action((args: object) => {
        const config = getConfig();
        const appKey = args['hide']
            ? '•'.repeat(config.application.key.length)
            : config.application.key;
        const clientKey = args['hide']
            ? '•'.repeat(config.client.key.length)
            : config.client.key;

        console.log(
`Soar Config
====================
Application Details

url: ${config.application.url || 'Not Set'}
key: ${appKey || 'Not Set'}
uses: 0

Client Details
url: ${config.client.url || 'Not Set'}
key: ${clientKey || 'Not Set'}

General
use strict mode:            ${config.logs.strictMode}
output debug messages:      ${config.logs.showDebug}
output http responses:      ${config.logs.showHttpLog}
output websocket responses: ${config.logs.showWsLog}

Error Directory: ${config.logs.errorOutDir || 'Not Set'}`
        );
    });

const main = new Command('config')
    .addHelpText('before', 'Manages the internal Soar configurations.')
    .addCommand(infoCmd)
    .addCommand(setupCmd);

export default main;
