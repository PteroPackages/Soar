import { Interface } from 'readline';
import { writeFileSync } from 'fs';
import { join } from 'path';
import log from '../log';

export async function writeFileResponse(name: string, data: string, writeLog: boolean): Promise<void> {
    return new Promise<void>(
        (res, rej) => {
            try {
                writeFileSync(
                    join(process.cwd(), name),
                    Buffer.from(data),
                    { encoding: 'utf-8' }
                );
                if (writeLog) log.success([
                    'saved request response at:',
                    join(process.cwd(), name)
                ]);
                res();
            } catch (err) {
                log.error(
                    'Internal Error',
                    [
                        'couldn\'t write response file',
                        err.message
                    ]
                );
                rej();
            }
        }
    );
}

async function prompt(reader: Interface, message: string): Promise<string> {
    reader.write(`[soar] ${message}\n[soar] >`);

    let out: string;
    await new Promise<void>(
        res => {
            reader.once('line', line => {
                out = line;
                res();
            });
        }
    );

    return Promise.resolve(out);
}

export async function getBoolInput(reader: Interface, message: string): Promise<boolean> {
    const out = await prompt(reader, message);
    if ('yes'.includes(out.slice(8))) return true;
    return false;
}

export async function getStringInput(reader: Interface, message: string, allowEmpty: boolean) {
    const out = await prompt(reader, message);
    if (!out.length && !allowEmpty) return getStringInput(reader, message, allowEmpty);
    return out;
}

function getMaxLength(data: object): number {
    let max = 0;

    for (const [k, v] of Object.entries(data)) {
        if (k.length > max) max = k.length;
        if (
            typeof v === 'object' &&
            !Array.isArray(v) &&
            v !== null
        ) {
            const t = getMaxLength(v);
            if (t > max) max = t;
        }
    }

    return max;
}

export function formatString(data: object): string {
    const max = getMaxLength(data);
    let res: string[] = [];

    if (Array.isArray(data)) return data.map(formatString).join('\n');

    for (let [key, val] of Object.entries(data)) {
        let fmt = `${key}: ${' '.repeat(max - key.length)}`;
        if (
            typeof val === 'object' &&
            val !== null
        ) {
            fmt += '<object ref>\n\n';
            fmt += formatString(val);
            fmt += '\n';
        } else {
            fmt += `${val}`;
        }
        res.push(fmt);
    }

    return res.join('\n');
}
