import yaml from 'yaml';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { error, fromError } from '../log';
import { parseStruct, Config, jsonStruct } from '../structs';

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

export function createConfig(options?: Config): void {
    if (!process.env.SOAR_PATH) error('MISSING_ENV', null, true);
    let config: Config;

    try {
        const temp = yaml.parse(readFileSync('../../config.ex.yml', 'utf-8'));
        config = parseStruct(temp);
    } catch (err) {
        fromError(err, true);
    }

    let data: object;
    if (options) {
        data = jsonStruct<Config>(compareConfigs(config, options));
    } else {
        data = jsonStruct(config);
    }

    try {
        writeFileSync(
            process.env.SOAR_PATH,
            yaml.stringify(data),
            { encoding: 'utf-8' }
        );
    } catch {
        error('MISSING_PERMISSIONS', null, true);
    }
}

export function updateConfig(newConfig: Config): void {}

function compareConfigs(_old: Config, _new: Config): Config {
    for (const [k, v] of Object.entries(_new)) {
        if (
            _old[k] === undefined ||
            _old[k] === ''
        ) _old[k] = v;
    }
    return _old;
}
