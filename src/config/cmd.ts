import { Command } from 'commander';
import { existsSync } from 'fs';
import { join } from 'path';
import log from '../log';
import { getBoolInput } from '../response';
import { getConfig, createConfig } from './funcs';

const infoCmd = new Command('info')
    .option('--local', 'Gets the local configuration for the workspace.', false)
    .option('-h, --hide', 'Hides the API keys from the command output.', false)
    .action((args: object) => {
        const config = getConfig(args['local']);
        const appKey = args['hide']
            ? '•'.repeat(config.application.key.length)
            : config.application.key;
        const clientKey = args['hide']
            ? '•'.repeat(config.client.key.length)
            : config.client.key;

        console.log(
`Soar ${args['local'] ? 'Local' : 'Global'} Config
====================
\x1b[4mApplication Details\x1b[0m
url: ${config.application.url || 'Not Set'}
key: ${appKey || 'Not Set'}

\x1b[4mClient Details\x1b[0m
url: ${config.client.url || 'Not Set'}
key: ${clientKey || 'Not Set'}

\x1b[4mGeneral\x1b[0m
use strict mode:            ${config.logs.strictMode}
output debug messages:      ${config.logs.showDebug}
output http responses:      ${config.logs.showHttpLog}
output websocket responses: ${config.logs.showWsLog}

Error Directory: ${config.logs.errorOutDir || 'Not Set'}`
        );
    });

const setupCmd = new Command('setup')
    .addHelpText('before', 'Setup a new Soar configuration.')
    .option('--local', 'Setup a local configuration for the workspace (default is global).', false)
    .option('--link [file]', 'Links the new config with another local config or the global config if not provided.')
    .option('-f, --force', 'Skips all confirmation prompts.', false)
    .action(async (args: object) => {
        const local: boolean = args['local'];
        const force: boolean = args['force'];
        let link: string | boolean = args['link'];
        let linkfp: string;

        if (local) {
            if (existsSync(join(process.cwd(), '.soar-local.yml'))) {
                log.notice('existing local config file found');
                if (!force) {
                    const res = getBoolInput('Do you want to overwrite this file? (y/n)');
                    if (!res) return;
                } else {
                    log.notice('overwrite mode forced for local config');
                }
            }
        }

        if (link) {
            if (typeof link === 'string') {
                if (!existsSync(link)) log.error(
                    'Not Found Error',
                    'The local config file path could not be resolved',
                    true
                );
                linkfp = link;
            } else {
                if (!existsSync(process.env.SOAR_PATH)) log.error('MISSING_CONFIG', null, true);
                linkfp = process.env.SOAR_PATH;
            }
        }

        const fp = local ?
            join(process.cwd(), '.soar-local.yml')
            : join(process.env.SOAR_PATH, 'config.yml');

        await createConfig(fp, linkfp);

        log.success([`setup a new ${local ? 'local' : 'global'} config at:`, fp]);
    });

export default [
    infoCmd,
    setupCmd
]
