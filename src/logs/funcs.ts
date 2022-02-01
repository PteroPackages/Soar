import { join } from 'path';
import {
    appendFileSync,
    existsSync,
    readFileSync,
    writeFileSync
} from 'fs';
import { version } from '../../package.json';
import { ReqLog } from '../structs';
import log from '../log';

export function fetchLogs() {
    if (!process.env.SOAR_PATH) log.error('MISSING_ENV', null, true);

    const fp = join(process.env.SOAR_PATH, 'logs/requests.log');
    if (!existsSync(fp)) make(fp);

    return parseLogs(readFileSync(fp, { encoding: 'utf-8' }));
}

function getLastRef() {
    try {
        const logs = readFileSync(
            join(process.env.SOAR_PATH, 'logs/requests.log'),
            { encoding: 'utf-8' }
        );
        return logs
            .split('\n')
            .filter(Boolean)
            .pop()
            .split(':')[0] || '0';
    } catch {
        return '0';
    }
}

export function createRequestLog(_log: ReqLog): void {
    if (!process.env.SOAR_PATH) log.error('MISSING_ENV', null, true);

    const fp = join(process.env.SOAR_PATH, 'logs/requests.log');
    if (!existsSync(fp)) make(fp);

    const fmt = `${_log.date}|${_log.method}|${_log.response}`+
        `|${_log.type}|${_log.domain.replace(/https?:\/\//g, '')}`+
        `|${_log.path}|${getLastRef()}\n`;

    try {
        appendFileSync(fp, fmt, { encoding: 'utf-8' });
    } catch {
        log.warn(`could not write log for '${Date.now()}'. Please check the application permissions.`);
    }
}

function make(path: string) {
    try {
        writeFileSync(path, `#${version}\n`, { encoding: 'utf-8' });
    } catch (err) {
        log.fromError(err, true);
    }
}

function parseLogs(data: string): ReqLog[] {
    const res: ReqLog[] = [];

    for (const line of data.split('\n')) {
        if (!line.length) continue;
        if (line.startsWith('#')) continue;
        const log = {} as ReqLog;
        const [d, m, s, t, b, p, r] = line.split('|');

        log.date = Number(d);
        log.method = m;
        log.response = Number(s);
        log.type = t;
        log.domain = b;
        log.path = p;
        log.ref = r;

        res.push(log);
    }

    return res;
}
