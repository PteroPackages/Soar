import { createInterface } from 'readline';
import * as writer from './writer';

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

export function handleCloseInterface(args: string[]) {}
