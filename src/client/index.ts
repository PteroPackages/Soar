import { Command } from 'commander';
import account from './account';

const root = new Command('client')
    .description('Commands for interacting with the Client API');

for (const cmd of account) root.addCommand(cmd);

export default root;
