import { Command } from 'commander';
import users from './users';
import servers from './servers';
import nodes from './nodes';
import locations from './locations';
import nests from './nests';

const root = new Command('app')
    .description('Commands for interacting with the Application API');

for (const cmd of users) root.addCommand(cmd);
for (const cmd of servers) root.addCommand(cmd);
for (const cmd of nodes) root.addCommand(cmd);
for (const cmd of locations) root.addCommand(cmd);
for (const cmd of nests) root.addCommand(cmd);

export default root;
