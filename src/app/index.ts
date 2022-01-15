import { Command } from 'commander';
import { handleRequest } from '../request';
import * as res from '../response';
import { parseUserGroup } from '../validate';
import * as log from '../log';
import Waiter from '../log/waiter';

const getUsersCmd = new Command('get-users')
    .addHelpText('before', 'Fetches all accounts from the panel (can specify or query with flags).')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .action(async (args: object) => {
        const options = parseUserGroup(args);
        let waiter: Waiter;

        if (!options.silent) {
            waiter = new Waiter(log.parse('%yfetching%R /application/users', 'info'))
                .onEnd(t => log.parse(`%gfetched%R /application/users (${t}ms taken)`, 'success'));
            waiter.start();
        }

        const data = await handleRequest('GET', '/api/application/users');
        if (!options.silent) {
            waiter.stop();
            log.print('Request Result:\n');
        }

        const out = res.handleCloseInterface(data, options);
        if (out) console.log(out);
    });

const main = new Command('app')
    .addCommand(getUsersCmd);

export default main;
