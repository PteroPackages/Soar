import yaml from 'yaml';
import { existsSync, readFileSync } from 'fs';
import { error } from '../log';
import { ERRORS } from '../errors';
import { Config } from '../structs';

export function getConfig() {
    if (!process.env.SOAR_PATH) error('PathError', ERRORS.MISSING_ENV, true);
    if (!existsSync(process.env.SOAR_PATH)) error('PathError', ERRORS.INVALID_ENV, true);

    try {
        const config = yaml.parse(readFileSync(process.env.SOAR_PATH, 'utf-8'));
        return config as Config;
    } catch {
        error('SoarError', ERRORS.CANNOT_READ_ENV, true);
    }
}

export function updateConfig(newConfig: Config): void {}

function compareConfigs(_old: Config, _new: Config)/*: Config */{}
