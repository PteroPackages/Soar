import { Command } from 'commander';
import { version } from '../package.json';

import app from './app';
import config from './config';

const root = new Command('soar')
    .version(version, '-v, --version')
    .addCommand(app)
    .addCommand(config);

root.parse(process.argv);
