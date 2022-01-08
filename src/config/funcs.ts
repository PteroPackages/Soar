import yaml from 'yaml';
import { existsSync, readFileSync } from 'fs';
import getError from '../errors';
import { Config } from '../structs';

export function getConfig() {
    if (!process.env.SOAR_PATH) getError('MISSING_ENV');
    if (!existsSync(process.env.SOAR_PATH)) getError('INVALID_ENV');

    try {
        const config = yaml.parse(readFileSync(process.env.SOAR_PATH, 'utf-8'));
        return config as Config;
    } catch {
        getError('CANNOT_READ_ENV');
    }
}

export function updateConfig(newConfig: Config): void {}

function compareConfigs(_old: Config, _new: Config)/*: Config */{}
