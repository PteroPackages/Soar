import { Command } from 'commander';
import Session from '../session';
import { buildServer, parseFlagOptions } from '../validate';
import log from '../log';

const getServersCmd = new Command('get-servers')
    .addHelpText('before', 'Fetches all servers from the panel.')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('--text', 'Send the response output as formatted text.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .option('--id <id>', 'The server ID to fetch.')
    .option('--uuid <uuid>', 'The UUID to query.')
    .option('--name <name>', 'The server name to query.')
    .option('--external <id>', 'The external server ID to query.')
    .option('--image <url>', 'The docker image URL to query.')
    .action(async (args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        const data = await session.handleRequest('GET', buildServer(args));
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request result:\n');
            console.log(out);
        }
    });

const suspendServerCmd = new Command('suspend')
    .addHelpText('before', 'Suspends a specified server.')
    .argument('<id>', 'The ID of the server.')
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        await session.handleRequest('POST', buildServer({ id, suspend: true }));
        if (!options.silent) log.success(`suspended server: ${id}`);
    });

    const unsuspendServerCmd = new Command('unsuspend')
    .addHelpText('before', 'Unsuspends a suspended server.')
    .argument('<id>', 'The ID of the server.')
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        await session.handleRequest('POST', buildServer({ id, suspend: true }));
        if (!options.silent) log.success(`unsuspended server: ${id}`);
    });

export default [
    getServersCmd,
    suspendServerCmd,
    unsuspendServerCmd
]
