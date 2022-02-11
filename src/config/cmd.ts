import { Command } from 'commander';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { join } from 'path';
import { createInterface } from 'readline';
import log from '../log';
import { getBoolInput } from '../session/response';
import { getConfig, createConfig } from './funcs';

const infoCmd = new Command('info')
    .option('--local', 'Gets the local configuration for the workspace.', false)
    .option('-h, --hide', 'Hides the API keys from the command output.', false)
    .action(async (args: object) => {
        const config = await getConfig(args['local']);
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
output debug messages:      ${config.logs.showDebug}
output http responses:      ${config.logs.showHttpLog}
output websocket responses: ${config.logs.showWsLog}
log http requests:          ${config.logs.logHttpRequests}
ignore warnings:            ${config.logs.ignoreWarnings}
output full response body:  ${config.logs.sendFullBody}`
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
        let mainfp = join(process.env.SOAR_PATH || '', 'config.yml');

        if (link) {
            if (typeof link === 'string') {
                if (!existsSync(link)) log.error(
                    'Not Found Error',
                    'the local config file path could not be resolved',
                    true
                );
                mainfp = link;
            } else {
                if (!existsSync(process.env.SOAR_PATH)) log.error('MISSING_CONFIG', null, true);
            }
        }

        if (local) {
            if (existsSync(join(process.cwd(), '.soar-local.yml'))) {
                log.info('existing local config file found');
                if (!force) {
                    const reader = createInterface(
                        process.stdin,
                        process.stdout
                    );

                    const res = await getBoolInput(reader, 'do you want to overwrite this file? (y/n)');
                    if (!res) return;
                } else {
                    log.info('overwrite mode forced for local config');
                }
            }

            if (!existsSync(mainfp)) log.error('MISSING_ENV', null, true);
            try {
                const linkData = readFileSync(mainfp, { encoding: 'utf-8' });
                writeFileSync(
                    join(process.cwd(), '.soar-local.yml'),
                    linkData, { encoding: 'utf-8' }
                );

                log.success([
                    'setup a new local config at:',
                    join(process.cwd(), '.soar-local.yml')
                ]);
            } catch (err) {
                log.fromError(err, true);
            }
        } else {
            await createConfig(link ? mainfp : null);
        }
    });

export default [
    infoCmd,
    setupCmd
]
