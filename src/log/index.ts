import errors from '../errors';

export const COLOURS = {
    RESET: '\x1b[0m',
    BOLD: '\x1b[1m',
    RED: '\x1b[31m',
    YELLOW: '\x1b[33m',
    GREEN: '\x1b[32m',
    CYAN: '\x1b[36m',
    BLUE: '\x1b[34m',
    MAGENTA: '\x1b[35m'
}

export const CODE_MAP = {
    '%R': COLOURS.RESET,
    '%B': COLOURS.BOLD,
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
        type?: 'info' | 'success' | 'notice' | 'debug' | 'warn' | 'error'
    ): string {
    for (const [k, c] of Object.entries<string>(CODE_MAP)) message = message.replaceAll(k, c);

    if (type) {
        switch (type) {
            case 'info': message = `${BASE} ${COLOURS.BLUE}info${COLOURS.RESET}: ${message}`; break;
            case 'success': message = `${BASE} ${COLOURS.GREEN}success${COLOURS.RESET}: ${message}`; break;
            case 'notice': message = `${BASE} ${COLOURS.BLUE}notice${COLOURS.RESET}: ${message}`; break;
            case 'debug': message = `${BASE} ${COLOURS.BOLD}debug${COLOURS.RESET} ${message}`; break;
            case 'warn': message = `${BASE} ${COLOURS.YELLOW}warning${COLOURS.RESET}: ${message}`; break;
            case 'error': message = `${BASE} ${COLOURS.RED}error${COLOURS.RESET}: ${message}`; break;
        }
    }

    return message;
}

export function print(message: string): void {
    console.log(BASE +' '+ parse(message));
}

export function info(message: string | string[]): void {
    const border = `${BASE} ${COLOURS.BLUE}info${COLOURS.RESET}: `;
    const fmt = Array.isArray(message) ? message : [message];
    console.log(fmt.map(m => border + m).join('\n'));
}

export function notice(message: string | string[]): void {
    const border = `${BASE} ${COLOURS.CYAN}notice${COLOURS.RESET}: `;
    const fmt = Array.isArray(message) ? message : [message];
    console.log(fmt.map(m => border + m).join('\n'));
}

export function debug(message: string | string[]): void {
    const border = `${BASE} ${COLOURS.BOLD}debug${COLOURS.RESET}: `;
    const fmt = Array.isArray(message) ? message : [message];
    console.log(fmt.map(m => border + m).join('\n'));
}

export function success(message: string | string[]): void {
    const border = `${BASE} ${COLOURS.GREEN}success${COLOURS.RESET}: `;
    const fmt = Array.isArray(message) ? message : [message];
    console.log(fmt.map(m => border + m).join('\n'));
}

export function warn(message: string | string[]): void {
    const border = `${BASE} ${COLOURS.YELLOW}warning${COLOURS.RESET}: `;
    const fmt = Array.isArray(message) ? message : [message];
    console.log(fmt.map(m => border + m).join('\n'));
}

export function error(name: string, message?: string | string[], exit?: boolean): void | never {
    const border = parse(`${BASE} %rerror%R: `);
    let fmt: string[] = [name];

    if (!message && errors.get(name)) {
        fmt.push(errors.get(name));
    } else {
        if (Array.isArray(message)) {
            fmt = fmt.concat(message);
        } else {
            if (message) fmt.push(message);
        }
    }

    console.log(fmt.map(m => border + m).join('\n'));
    if (exit) process.exit(1);
}

export function fromError(_error: Error, exit?: boolean): void | never {
    return error('Internal Error', _error.stack.split('\n'), exit);
}

interface pteroError {
    errors:{
        code:   string;
        status: string;
        detail: string;
    }[];
}

export function fromPtero(data: pteroError, exit?: boolean): void | never {
    error(
        'API Request Error',
        `Pterodactyl panel returned ${data.errors.length} error${data.errors.length > 1 ? 's' : ''}.`
    );

    for (const err of data.errors) {
        error('');
        error(err.code, err.detail || '[no details received]');
    }

    if (data.errors.some(e => e.status === '403')) {
        notice(
            'Please ensure that your API key has the necessary'+
            ' read/write permissions before making requests.'
        );
    }

    if (exit) process.exit(1);
}

export default {
    parse,
    print,
    info,
    notice,
    debug,
    success,
    warn,
    error,
    fromError,
    fromPtero
}
