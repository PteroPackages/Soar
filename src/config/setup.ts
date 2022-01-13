import { Command } from 'commander';
import { existsSync } from 'fs';
import * as log from '../log';
import { getBoolInput } from '../response';
import { createConfig } from './funcs';

const cmd = new Command('setup')
    .option('-f, --force', 'Skips the confirmation prompt.', false)
    .action((args: object) => {
        const shouldSkip: boolean = args['force'];

        if (!process.env.SOAR_PATH) log.error('MISSING_ENV', null, true);
        if (!existsSync(process.env.SOAR_PATH) && !shouldSkip) {
            log.info('Soar directories not found.');
            const res = getBoolInput(`Should new directories be set at '${process.env.SOAR_PATH}'?`);
            if (!res) process.exit(0);
        }

        createConfig();
        log.success('Successfully setup Soar directories!');
    });

export default cmd;
