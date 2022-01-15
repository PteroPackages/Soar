import { Command } from 'commander';

import users from './users';

const root = new Command('app');
for (const cmd of users) root.addCommand(cmd);

export default root;
