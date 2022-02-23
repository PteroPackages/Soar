import { createInterface } from 'readline';
import { Command, Option } from 'commander';
import Session from '../session';
import { buildApiKey, parseFlagOptions } from '../validate';
import log from '../log';
import { getBoolInput } from '../session/response';
import { getConfig } from '../config/funcs';

const getKeysCmd = new Command('get-keys')
    .description('Fetches the client account api keys')
    .addHelpText('before', 'Fetches the client account api keys')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('client', options);

        const data = await session.handleRequest('GET', buildApiKey());
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request result:\n');
            console.log(out);
        }
    });

const deleteKeyCmd = new Command('delete-key')
    .description('Deletes a specified key from the client account')
    .addHelpText('before', 'Deletes a specified key from the client account')
    .argument('<id>', 'The ID of the key to delete')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-f, --force', 'Force delete the key', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);
        const config = await getConfig(true);

        if (config.client.key.startsWith(id) && !args['force']) {
            const reader = createInterface(
                process.stdin,
                process.stdout
            );

            log.warn([
                'this api key is being used for soar requests',
                'deleting it will prevent future commands from working'
            ]);
            const res = await getBoolInput(reader, 'do you want to continue? (y/n)');
            if (!res) return;
        }

        await new Session('client', options)
            .handleRequest('DELETE', buildApiKey(id));
        if (!options.silent) log.success(`deleted api key: ${id}`);
    });

export default [
    getKeysCmd,
    deleteKeyCmd
]
