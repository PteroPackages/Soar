import { Command, Option } from 'commander';
import Session from '../session';
import { parseFlagOptions, buildNest } from '../validate';
import log from '../log';

const getNestsCmd = new Command('get-nests')
    .description('Fetches nests from the panel')
    .addHelpText('before', 'Fetches all nests from the panel (can specify with flags)')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .option('--id <id>', 'The node location ID to fetch')
    .option('--eggs', 'Include eggs in the request', false)
    .option('--servers', 'Include servers in the request', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        const data = await session.handleRequest('GET', buildNest(args));
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request results:\n');
            console.log(out);
        }
    });

export default [
    getNestsCmd
]
