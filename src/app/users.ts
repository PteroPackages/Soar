import { Command } from 'commander';
import Session from '../session';
import parseDiffView, { highlight } from '../session/view';
import { buildUser, parseFlagOptions } from '../validate';
import log from '../log';

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
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        const data = await session.handleRequest('GET', buildUser(args));
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request result:\n');
            console.log(out);
        }
    });

const createUserCmd = new Command('create-user')
    .addHelpText('before', 'Creates a new user account on the panel.')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('--text', 'Send the response output as formatted text.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .option('-d, --data <json>', 'The json data to create the user with.')
    .action(async (args: object) => {
        const options = parseFlagOptions(args);

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

        const missing: string[] = [];
        for (const key of ['username', 'email', 'first_name', 'last_name', 'language']) {
            if (key in json) continue;
            missing.push(key);
        }
        if (missing.length) log.error(
            'Argument Error',
            [
                `Missing required key${missing.length > 1 ? 's' : ''}:`,
                missing.join(', ')
            ],
            true
        );

        const session = new Session('application', options);
        await session.handleRequest('POST', buildUser({}), json);
        const data = await session.handleRequest('GET', buildUser({ email: json['email'] }));
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('account created! request result:\n');
            console.log(out);
        }
    });

const updateUserCmd = new Command('update-user')
    .addHelpText('before', 'Updates a specified user account.')
    .argument('<id>', 'The ID of the user account to update.')
    .option('--json', 'Send the response output as JSON.', false)
    .option('--yaml', 'Send the response output as YAML.', true)
    .option('--text', 'Send the response output as formatted text.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .option('-d, --data <json>', 'The json data to update the user with.')
    .option('--no-diff', 'Don\'t show the properties changed in the request.', false)
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);

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

        const session = new Session('application', options);
        const user = await session.handleRequest('GET', buildUser({ id }));
        if (!user) log.error('NOT_FOUND_USER', null, true);

        json['username'] ||= user['attributes']['username'];
        json['email'] ||= user['attributes']['email'];
        json['first_name'] ||= user['attributes']['first_name'];
        json['last_name'] ||= user['attributes']['last_name'];
        json['language'] ||= user['attributes']['language'];
        json['password'] ||= null;

        const data = await session.handleRequest('PATCH', buildUser({ id }), json);
        const out = await session.handleClose(data, options);

        if (out && args['diff']) {
            const view = parseDiffView(options.responseType, user, data);
            log.success(log.parse(
                `made %c${view.totalChanges}%R changes`+
                ` (%g+${view.additions}%R | %r-${view.subtractions}%R)`
            ));
            console.log('\n'+ highlight(view.output));
        } else {
            log.success(`updated user account: ${id}`);
        }
    });

const deleteUserCmd = new Command('delete-user')
    .addHelpText('before', 'Deletes a specified user account.')
    .argument('<id>', 'The ID of the user account to delete.')
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        await session.handleRequest('DELETE', buildUser({ id }));
        if (!options.silent) log.success(`deleted user account: ${id}`);
    });

export default [
    getUsersCmd,
    createUserCmd,
    updateUserCmd,
    deleteUserCmd
]
