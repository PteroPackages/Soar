import { Command } from 'commander';

import commands from './cmd';

const root = new Command('logs')
    .addHelpText('before', 'Manages Soar request logs.');

for (const cmd of commands) root.addCommand(cmd);

export default root;
