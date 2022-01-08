import { Command } from 'commander';
import { handleRequest } from '../request';
import { parseStruct, AppUser } from '../structs';

const getUsersCmd = new Command('get-users')
    .addHelpText('before', 'Fetches all accounts from the panel (can specify or query with flags).')
    .action(async () => {
        const data = await handleRequest('GET', '/api/application/users');
        const user = parseStruct<AppUser>(data);

        console.log(user);
    });

const main = new Command('app')
    .addCommand(getUsersCmd);

export default main;
