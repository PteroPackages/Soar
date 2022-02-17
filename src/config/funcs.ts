import yaml from 'yaml';
import { join } from 'path';
import { exec, ExecException } from 'child_process';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { Config } from '../structs';
import { parseConfig } from '../validate';
import log from '../log';

function run(cmd: string): Promise<[string, string | ExecException]> {
    let _res: [string, string | ExecException];

    return new Promise<[string, string | ExecException]>(
        res => {
            exec(cmd, (err: ExecException, out: string) => {
                if (err) _res = [null, err];
                else _res = [out, null];
            }).on('close', () => res(_res));
        }
    );
}

export async function getConfig(checkLocal: boolean = false): Promise<Config> {
    let fp: string;

    if (
        checkLocal &&
        existsSync(join(process.cwd(), '.soar-local.yml'))
    ) fp = join(process.cwd(), '.soar-local.yml');

    if (!fp) {
        if (!process.env.SOAR_PATH) log.error('MISSING_ENV', null, true);
        fp = join(process.env.SOAR_PATH, 'config.yml');
        if (!existsSync(fp)) log.error('INVALID_ENV', null, true);
    }

    try {
        const config = yaml.parse(readFileSync(fp, 'utf-8'));
        const err = parseConfig(config);
        if (err) log.error(
            'Config Error',
            [
                err,
                'make sure to update the necessary config option at:',
                fp
            ],
            true
        );
        return new Promise<Config>(res => res(config));
    } catch {
        log.error('CANNOT_READ_ENV', null, true);
    }
}

export function getConfigKey(config: Config, key: string): string[] {
    const [base, main] = key.split('.');
    if (main in config[base]) return [base, main];
    return [];
}

export async function createConfig(link?: string) {
    const lib = process.platform === 'win32'
            ? 'C:\\soar\\'
            : '/soar/';

    if (!process.env.SOAR_PATH) {
        log.info('soar library not found, attempting to fetch directly...');
        let [res, err] = await run('git --version');
        if (err) {
            err = err as ExecException;
            log.error(
                'Exec Error',
                err.message.includes('not found') || err.message.includes('not recognised')
                    ? 'git cli is required to continue'
                    : (err as ExecException).message,
                true
            );
        }

        if (existsSync(`${lib}bin`)) {
            log.info('existing soar library found, attempting clean...');
            [res, err] = await run(
                lib.includes('C:')
                    ? 'rmdir /S /Q C:\\soar\\bin'
                    : 'rm -rf /soar/bin'
            );
            if (err) {
                err = err as ExecException;
                log.error(
                    'Internal Error',
                    `could not remove existing soar library files at: '${lib}bin'`,
                    true
                );
            }
        }

        [res, err] = await run(`git clone https://github.com/PteroPackages/soar-ts.git ${lib}bin`);
        if (err) {
            err = err as ExecException;
            log.error(
                'Internal Error',
                [
                    ...err.message.split('\n'),
                    `source: ${err.cmd}`,
                    `code: ${err.code}`
                ],
                true
            );
        }

        log.success(`cloned soar library into ${lib}bin`);
        log.warn([
            `please set the environment variable 'SOAR_PATH' to ${lib.slice(0, -1)}`,
            `command: '%bset SOAR_PATH=${lib.slice(0, -1)}%R'`
        ]);
        process.env.SOAR_PATH = lib.slice(0, -1);
    }

    const tmpl = readFileSync(
        link || join(process.env.SOAR_PATH, 'bin/config.ex.yml'),
        'utf-8'
    );
    try {
        writeFileSync(`${lib}config.yml`, tmpl, { encoding: 'utf-8' });
    } catch (err) {
        log.error(
            'Internal Error',
            err.message.includes('permission') || err.message.includes('denied')
                ? 'missing the required read/write permissions to continue'
                : err.message,
            true
        );
    }
}

export function updateConfig(
    config: Config,
    key: string[],
    value: string,
    local: boolean
): void {
    if (['application', 'client'].includes(key[0])) {
        config[key[0]][key[1]] = value;
    } else {
        switch (value.toLowerCase()) {
            case 'true': config[key[0]][key[1]] = true; break;
            case 'false': config[key[0]][key[1]] = false; break;
            default: log.error(
                'Argument Error',
                `invalid config option value '${typeof value}'`,
                true
            );
        }
    }

    const err = parseConfig(config);
    if (err) log.error('Config Error', err, true);

    const fp = local
        ? join(process.cwd(), '.soar-local.yml')
        : process.env.SOAR_PATH;

    try {
        writeFileSync(fp, yaml.stringify(config), { encoding: 'utf-8' });
    } catch (err) {
        log.fromError(err, true);
    }
}
