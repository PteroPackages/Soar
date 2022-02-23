import { Command } from 'commander';
import account from './account';
import apikeys from './apikeys';
import twoFactor from './twoFactor';

const root = new Command('client')
    .description('Commands for interacting with the Client API');

for (const cmd of account) root.addCommand(cmd);
for (const cmd of apikeys) root.addCommand(cmd);
for (const cmd of twoFactor) root.addCommand(cmd);

export default root;
