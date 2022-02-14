import { Command } from 'commander';
import commands from './cmd';

const root = new Command('config')
    .description('Manages the global and local Soar config');

for (const cmd of commands) root.addCommand(cmd);

export default root;
