import { Command } from 'commander';
import { version } from '../package.json';
import { fromError } from './log';

import app from './app';
import client from './client';
import config from './config';
import log from './logs';

const root = new Command('soar')
    .version(`v${version}-beta (build: unknown)`, '-v, --version')
    .addCommand(app)
    .addCommand(client)
    .addCommand(config)
    .addCommand(log);

try {
    root.parse(process.argv);
} catch (err) {
    fromError(err, true);
}
