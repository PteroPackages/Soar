import { Command } from 'commander';
import { handleRequest } from '../request';
import * as res from '../response';
import { buildUser, parseUserGroup } from '../validate';
import * as log from '../log';
import Waiter from '../log/waiter';

const getUsersCmd = new Command('get-users')
    .addHelpText('before', 'Fetches all accounts from the panel (can specify or query with flags).')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .option('--id <userID>', 'The user ID to fetch.')
    .option('--email <email>', 'The email to query.')
    .option('--username <name>', 'The user name to query.')
    .option('--uuid <uuid>', 'The UUID to query.')
    .option('--external <id>', 'The external user ID to query.')
    .action(async (args: object) => {
        const options = parseUserGroup(args);
        let waiter: Waiter;

        if (!options.silent) {
            waiter = new Waiter(log.parse('%yfetching%R /application/users', 'info'))
                .onEnd(t => log.parse(`%gfetched%R /application/users (${t}ms taken)`, 'info'))
                .start();
        }

        const data = await handleRequest('GET', buildUser(args));
        if (!options.silent) {
            waiter.stop();
            log.info('request result:\n');
        }

        const out = res.handleCloseInterface(data, options);
        if (out) console.log(out);
    });

const main = new Command('app')
    .addCommand(getUsersCmd);

export default main;
