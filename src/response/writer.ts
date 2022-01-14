import { writeFileSync } from 'fs';
import { join } from 'path';
import yaml from 'yaml';

export function formatString(res: object): string {
    return '';
}

export function formatYAML(res: object): string {
    return yaml.stringify(res);
}

export function writeFileResponse(name: string, res: string) {
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
