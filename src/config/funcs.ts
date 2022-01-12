import yaml from 'yaml';
import { existsSync, readFileSync } from 'fs';
import { error } from '../log';
import { parseStruct, Config } from '../structs';

export function getConfig() {
    if (!process.env.SOAR_PATH) error('MISSING_ENV', null, true);
    if (!existsSync(process.env.SOAR_PATH)) error('INVALID_ENV', null, true);

    try {
        const config = yaml.parse(readFileSync(process.env.SOAR_PATH, 'utf-8'));
        return parseStruct<Config>(config);
    } catch {
        error('CANNOT_READ_ENV', null, true);
    }
}

export function updateConfig(newConfig: Config): void {}

function compareConfigs(_old: Config, _new: Config)/*: Config */{}
