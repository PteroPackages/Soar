import { Command } from 'commander';
import { version } from '../package.json';
import { fromError } from './log';

import app from './app';
import config from './config';

const root = new Command('soar')
    .version(version, '-v, --version')
    .addCommand(app)
    .addCommand(config);

try {
    root.parse(process.argv);
} catch (err) {
    fromError(err, true);
}
