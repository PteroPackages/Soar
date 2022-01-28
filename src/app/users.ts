import { Command } from 'commander';
import { handleRequest } from '../request';
import * as res from '../response';
import parseDiffView, { highlight } from '../response/view';
import { buildUser, parseUserGroup } from '../validate';
import * as log from '../log';
import Waiter from '../log/waiter';

const getUsersCmd = new Command('get-users')
    .addHelpText('before', 'Fetches all accounts from the panel (can specify or query with flags).')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('--text', 'Send the response output as formatted text.', false)
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

const updateUsersCmd = new Command('update-users')
    .addHelpText('before', 'Updates a specified user account.')
    .argument('<id>', 'The ID of the user to update.')
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .option('-d, --data <json>', 'The json data to update the user with.')
    .option('-c, --changes', 'Shows the properties changed in the request.', false)
    .action(async (id: string, args: object) => {
        const options = parseUserGroup(args);
        let waiter: Waiter;

        if (!options.silent) {
            waiter = new Waiter(log.parse('%yfetching%R /application/users', 'info'))
                .onEnd(t => log.parse(`%gfetched%R /application/users (${t}ms taken)`, 'info'))
                .start();
        }

        let json: object;
        try {
            json = JSON.parse(args['data']);
        } catch (err) {
            log.error(
                'Argument Error',
                [
                    'Couldn\'t parse JSON data argument:',
                    err.message
                ],
                true
            );
        }
        if (!Object.entries(json).length) log.error(
            'Argument Error',
            'No JSON was provided to update.',
            true
        );

        const user = await handleRequest('GET', buildUser({ id }));
        if (!user) log.error('NOT_FOUND_USER', null, true);

        json['username'] ||= user['attributes']['username'];
        json['email'] ||= user['attributes']['email'];
        json['first_name'] ||= user['attributes']['first_name'];
        json['last_name'] ||= user['attributes']['last_name'];
        json['language'] || user['attributes']['language'];
        json['password'] ||= null;

        const data = await handleRequest('PATCH', buildUser({ id }), json);
        if (!options.silent) waiter.stop();

        const out = res.handleCloseInterface(data, options);
        if (out) {
            const view = parseDiffView('yaml', user, data);
            log.info(log.parse(
                `made %c${view.totalChanges}%R changes (%g+${view.additions}%R | %r-${view.subtractions}%R)`
            ));
            console.log('\n'+ highlight(view.output));
        }
    });

export default [
    getUsersCmd,
    updateUsersCmd
]
