import { Command } from 'commander';

import commands from './cmd';

const root = new Command('config')
    .addHelpText('before', 'Manages the internal Soar configurations.');

for (const cmd of commands) root.addCommand(cmd);

export default root;
