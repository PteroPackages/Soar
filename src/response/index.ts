import { createInterface } from 'readline';
import { writeFileSync } from 'fs';
import { join } from 'path';
import yaml from 'yaml';
import { FlagOptions } from '../structs';

export function getBoolInput(message: string): boolean {
    const input = createInterface(process.stdin);
    let res: boolean;
    input.question(message, ans => {
        if ('yes'.includes(ans.toLowerCase())) res = true;
        else res = false;
    });
    return res;
}

export function getStringInput(message: string, allowEmpty: boolean): string {
    const input = createInterface(process.stdin);
    let res: string;
    input.question(message, ans => res = ans);
    if (!res.length) return getStringInput(message, allowEmpty);
    return res;
}

export function getOptionInput(
        message: string,
        options: string[],
        _default?: string,
        errorMessage?: string
    ): string {
    const res = getStringInput(message, _default !== null);
    if (!res) return _default;
    if (!options.includes(res)) {
        console.log(errorMessage || `Invalid option '${res.length > 10 ? res.slice(0, 10)+'...' : res}'`);
        return getOptionInput(message, options, _default);
    }
    return res;
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

function formatString(data: object): string {
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

function writeFileResponse(name: string, res: string) {
    try {
        writeFileSync(
            join(process.cwd(), name),
            Buffer.from(res),
            { encoding: 'utf-8' }
        );
    } catch (err) {
        console.error(err);
    }
}

export function handleCloseInterface(data: object, options: FlagOptions): string | void {
    let parsed: string;

    switch (options.responseType) {
        case 'text': parsed = formatString(data); break;
        case 'yaml': parsed = yaml.stringify(data); break;
        default: parsed = JSON.stringify(data); break;
    }

    if (options.writeFile.length) {
        writeFileResponse(options.writeFile, parsed);
    } else if (options.prompt && !options.silent) {
        const res = getBoolInput('Should this request be saved? (y/n)');
        if (res) {
            const fp = getStringInput('Enter the file path to save to: ', true);
            if (fp) writeFileResponse(`soar_log_${Date.now()}`, parsed);
        }
    }

    return parsed;
}
