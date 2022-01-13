import { Command } from 'commander';
import { handleRequest } from '../request';
import * as res from '../response';
import { parseUserGroup } from '../validate';

const getUsersCmd = new Command('get-users')
    .addHelpText('before', 'Fetches all accounts from the panel (can specify or query with flags).')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .action(async (args: object) => {
        const options = parseUserGroup(args);
        const data = await handleRequest('GET', '/api/application/users');
        const out = res.handleCloseInterface(data, options);
        if (out) console.log(out);
    });

const main = new Command('app')
    .addCommand(getUsersCmd);

export default main;
