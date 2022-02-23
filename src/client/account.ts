import { Command, Option } from 'commander';
import Session from '../session';
import parseDiffView, { highlight } from '../session/view';
import { parseFlagOptions } from '../validate';
import log from '../log';

const getAccountCmd = new Command('get-account')
    .description('Fetches the client account from the panel')
    .addHelpText('before', 'Fetches the client account from the panel (can specify or query with flags)')
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

        const data = await session.handleRequest('GET', '/api/client/account');
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request result:\n');
            console.log(out);
        }
    });

const updateAccountCmd = new Command('update-account')
    .description('Updates the client account on the panel')
    .addHelpText('before', 'Updates the client account on the panel with flags')
    .option('--json', 'Send the response output as JSON', false)
    .option('--yaml', 'Send the response output as YAML', true)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .option('--email <email>', 'The new email address (requires password)')
    .option('--new <password>', 'The new password (required for updating)')
    .requiredOption('--pass <password>', 'The account password (required for updating)')
    .option('--no-diff', 'Don\'t show the properties changed in the request', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (args: object) => {
        let path = '/api/client/account/';
        let json: object;

        if (args['email']) {
            json = { email: args['email'], password: args['pass'] };
            path += 'email';
        } else {
            if (!args['new']) log.error(
                'Argument Error',
                "'--new' and '--pass' are required to update password",
                true
            );
            json = {
                current_password: args['pass'],
                password: args['new'],
                password_confirmation: args['new']
            };
            path += 'password';
        }

        const options = parseFlagOptions(args);
        const session = new Session('client', options);

        const user = await session.handleRequest('GET', '/api/client/account');
        await session.handleRequest('PUT', path, json);
        const data = await session.handleRequest('GET', '/api/client/account');
        const out = await session.handleClose(data, options);

        if (out && args['diff'] && !args['new']) {
            const view = parseDiffView(options.responseType, user, data);
            log.print(
                'success',
                `made %c${view.totalChanges}%R changes`+
                ` (%g+${view.additions}%R | %r-${view.subtractions}%R)`,
                false
            );
            console.log(
                '\n'+ (session.config.logs.useColour
                ? highlight(view.output)
                : view.output)
            );
        } else {
            log.success('updated client account');
        }
    });

export default [
    getAccountCmd,
    updateAccountCmd
]
