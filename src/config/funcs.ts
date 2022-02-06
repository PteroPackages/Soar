import yaml from 'yaml';
import { join } from 'path';
import { exec, ExecException } from 'child_process';
import log from '../log';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { parseStruct, Config } from '../structs';

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
        if (!process.env.SOAR_PATH) log.error('MISSING_CONFIG', null, true);
        fp = join(process.env.SOAR_PATH, 'config.yml');
        if (!existsSync(fp)) log.error('INVALID_ENV', null, true);
    }

    try {
        const config = yaml.parse(readFileSync(fp, 'utf-8'));
        return new Promise<Config>(res => res(parseStruct<Config>(config)));
    } catch {
        log.error('CANNOT_READ_ENV', null, true);
    }
}

export async function createConfig(path: string, link?: string) {
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

        const lib = process.platform === 'win32'
            ? 'C:\\soar\\'
            : '/soar/';

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
            `command: 'set SOAR_PATH=${lib.slice(0, -1)}'`
        ]);
        process.env.SOAR_PATH = lib.slice(0, -1);
    }

    const tmpl = readFileSync(
        link || join(process.env.SOAR_PATH, 'bin/config.ex.yml'),
        'utf-8'
    );
    try {
        writeFileSync(path, tmpl, { encoding: 'utf-8' });
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
