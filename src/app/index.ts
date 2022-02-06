import { Command } from 'commander';

import users from './users';
import servers from './servers';
import nodes from './nodes';

const root = new Command('app');
for (const cmd of users) root.addCommand(cmd);
for (const cmd of servers) root.addCommand(cmd);
for (const cmd of nodes) root.addCommand(cmd);

export default root;
