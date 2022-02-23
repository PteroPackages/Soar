import { Command, Option } from 'commander';
import Session from '../session';
import { parseFlagOptions } from '../validate';
import log from '../log';

const get2FAInfo = new Command('get-2fa')
    .description('Fetches the 2fa totp codes for the account')
    .addHelpText('before', 'Fetches the 2fa totp codes for the account')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .option('--saved', 'Gets the saved codes, if any', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (args: object) => {
        // TODO
        // if (args['saved']) {}

        const options = parseFlagOptions(args);
        const session = new Session('client', options);

        const data = await session.handleRequest('GET', '/api/client/account/two-factor');
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request result:\n');
            console.log(out);
        }
    });

const enable2FACmd = new Command('enable-2fa')
    .description('Fetches the 2fa totp codes for the account')
    .addHelpText('before', 'Fetches the 2fa totp codes for the account')
    .argument('<code>', 'The totp 2fa code to enable with')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (code: string, args: object) => {
        log.warn(
            'this feature has been known to be problematic across all PteroPackages libraries'
        );

        const options = parseFlagOptions(args);
        const session = new Session('client', options);

        const data = await session.handleRequest(
            'POST',
            '/api/client/account/two-factor',
            { code: args['code'] }
        );
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request result:\n');
            console.log(out);
        }
    });

const disable2FACmd = new Command('disable-2fa')
    .description('Fetches the 2fa totp codes for the account')
    .addHelpText('before', 'Fetches the 2fa totp codes for the account')
    .argument('<code>', 'The totp 2fa code to enable with')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (password: string, args: object) => {
        const options = parseFlagOptions(args);
        await new Session('client', options)
            .handleRequest(
                'DELETE',
                '/api/client/account/two-factor',
                { password }
            );

        if (!options.silent) log.success('disabled 2fa');
    });

export default [
    get2FAInfo,
    enable2FACmd,
    disable2FACmd
]
