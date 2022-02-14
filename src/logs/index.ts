import { Command } from 'commander';
import commands from './cmd';

const root = new Command('logs')
    .description('Commands for interacting with the logging system');

for (const cmd of commands) root.addCommand(cmd);

export default root;
