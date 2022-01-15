import { Command } from 'commander';
import { handleRequest } from '../request';
import { handleCloseInterface } from '../response';
import { buildServer, parseServerGroup } from '../validate';
import * as log from '../log';
import Waiter from '../log/waiter';

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
        const options = parseServerGroup(args);
        let waiter: Waiter;

        if (!options.silent) {
            waiter = new Waiter(log.parse('%yfetching%R /application/servers', 'info'))
                .onEnd(t => log.parse(`%gfetched%R /application/servers (${t}ms taken)`, 'info'))
                .start();
        }

        const data = await handleRequest('GET', buildServer(args));
        if (!options.silent) {
            waiter.stop();
            log.info('request result:\n');
        }

        const out = handleCloseInterface(data, options);
        if (out) console.log(out);
    });

export default [
    getServersCmd
]
