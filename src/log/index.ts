import * as errors from '../errors';

export const COLOURS = {
    BASE: '\x1b[',
    RESET: '\x1b[0m',
    RED: '\x1b[31m',
    YELLOW: '\x1b[33m',
    GREEN: '\x1b[32m',
    CYAN: '\x1b[36m',
    BLUE: '\x1b[34m',
    MAGENTA: '\x1b[35m'
}

export const CODE_MAP = {
    '%R': COLOURS.RESET,
    '%r': COLOURS.RED,
    '%y': COLOURS.YELLOW,
    '%g': COLOURS.GREEN,
    '%c': COLOURS.CYAN,
    '%b': COLOURS.BLUE,
    '%m': COLOURS.MAGENTA
}

const BASE = `[${COLOURS.CYAN}soar${COLOURS.RESET}]`;

export function parse(
        message: string,
        type?: 'info' | 'success' | 'notice' | 'warn' | 'error'
    ): string {
    for (const [k, c] of Object.entries<string>(CODE_MAP)) message = message.replaceAll(k, c);

    if (type) {
        switch (type) {
            case 'info': message = `${BASE} ${COLOURS.BLUE}info${COLOURS.RESET}: ${message}`; break;
            case 'success': message = `${BASE} ${COLOURS.GREEN}success${COLOURS.RESET}: ${message}`; break;
            case 'notice': message = `${BASE} ${COLOURS.BLUE}notice${COLOURS.RESET}: ${message}`; break;
            case 'warn': message = `${BASE} ${COLOURS.YELLOW}warning${COLOURS.RESET}: ${message}`; break;
            case 'error': message = `${BASE} ${COLOURS.RED}error${COLOURS.RESET}: ${message}`; break;
        }
    }

    return message;
}

export function print(message: string): void {
    console.log(parse(message));
}

export function info(message: string): void {
    console.log(`${BASE} ${COLOURS.BLUE}info${COLOURS.RESET}: ${message}`);
}

export function success(message: string): void {
    console.log(`${BASE} ${COLOURS.GREEN}success${COLOURS.RESET}: ${message}`);
}

export function warn(message: string | string[]): void {
    const border = parse(`${BASE} %ywarning%R: `);
    const fmt = Array.isArray(message) ? message : [message];

    console.log(fmt.map(m => border + m).join('\n'));
}

export function error(name: string, message?: string | string[], exit?: boolean): void | never {
    const border = parse(`${BASE} %rerror%R: `);
    let fmt: string[] = [name];

    if (!message) {
        if (errors.tryGet(name)) {
            fmt.push(errors.get(name));
        } else {
            fmt.push('INVALID LOG MESSAGE');
        }
    } else {
        if (Array.isArray(message)) {
            fmt = fmt.concat(message);
        } else {
            fmt.push(message);
        }
    }

    console.log(fmt.map(m => border + m).join('\n'));
    if (exit) process.exit(1);
}

export function fromError(_error: Error, exit?: boolean): void | never {
    return error('Internal Errror', _error.stack, exit);
}
